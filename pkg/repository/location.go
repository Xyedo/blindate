package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
)

type Location interface {
	InsertNewLocation(location *entity.Location) (int64, error)
	UpdateLocation(location *entity.Location) (int64, error)
	GetLocationByUserId(id string) (*entity.Location, error)
	GetClosestUser(userId, geom string, limit int) ([]domain.BigUser, error)
}

func NewLocation(db *sqlx.DB) *location {
	return &location{
		conn: db,
	}
}

type location struct {
	conn *sqlx.DB
}

func (l *location) InsertNewLocation(location *entity.Location) (int64, error) {
	query := `
		INSERT INTO locations(user_id, geog, created_at, updated_at)
		VALUES($1, ST_GeomFromText($2), $3, $3)`
	now := time.Now()
	args := []any{location.UserId, location.Geog, now}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := l.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	ret, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}
	location.CreatedAt = now
	location.UpdatedAt = now
	return ret, nil
}

func (l *location) UpdateLocation(location *entity.Location) (int64, error) {
	query := `
		UPDATE locations SET geog = ST_GeomFromText($1), updated_at = $2
		WHERE user_id = $3`
	args := []any{location.Geog, time.Now(), location.UserId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := l.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, nil
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return rows, nil
}

func (l *location) GetLocationByUserId(id string) (*entity.Location, error) {
	query := `
		SELECT 
			user_id,
			ST_AsText(geog) as geog,
			created_at, 
			updated_at 
		FROM locations 
		WHERE user_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var location entity.Location
	err := l.conn.GetContext(ctx, &location, query, id)
	if err != nil {
		return nil, err
	}
	return &location, nil

}

func (l *location) GetClosestUser(userId, geom string, limit int) ([]domain.BigUser, error) {
	query := `
		SELECT 
			u.id as user_id,
			u.alias,
			u.dob,
			b.gender,
			b.from_loc,
			b.height,
			b.education_level,
			b.drinking,
			b.smoking,
			b.relationship_pref,
			b.looking_for,
			b.zodiac,
			b.kids,
			b.work,
			i.id as bio_id,
			i.bio,
			ARRAY(
				SELECT hobbie 
				FROM hobbies
				WHERE i.id IS NOT NULL AND interest_id = i.id
			) as interest_hobbies,
			ARRAY(
				SELECT movie_serie 
				FROM movie_series
				WHERE i.id IS NOT NULL AND interest_id = i.id
			) as interest_movie_series,
			ARRAY(
				SELECT travel 
				FROM traveling
				WHERE i.id IS NOT NULL AND interest_id = i.id
			) as interest_traveling,
			ARRAY(
				SELECT sport 
				FROM sports
				WHERE i.id IS NOT NULL AND interest_id = i.id
			) as interest_sport
		FROM locations l
		JOIN users u
			ON u.id = l.user_id
		LEFT JOIN basic_info b
			ON b.user_id = u.id
		LEFT JOIN interests i
			ON i.user_id = u.id
		WHERE NOT EXISTS (
			SELECT 1
			FROM match m
			WHERE 
				m.request_to = u.id OR
				m.request_from = u.id
		) AND u.id != $3
		ORDER BY l.geog <-> ST_GeomFromText($1)
		LIMIT $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matchs := make([]domain.BigUser, 0)
	rows, err := l.conn.QueryxContext(ctx, query, geom, limit, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}(rows)
	for rows.Next() {
		bigUser, err := l.createBigUser(rows)
		if err != nil {
			return nil, err
		}
		matchs = append(matchs, bigUser)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return matchs, nil
}
func (*location) createBigUser(row sqlx.ColScanner) (domain.BigUser, error) {
	var newBigUser domain.BigUser
	var newBasicInfo entity.BasicInfo
	var newBasicInfoGender sql.NullString
	var newBasicInfoLookingFor sql.NullString
	var newMatchBioId sql.NullString
	var newMatchBio sql.NullString
	var hobbies pq.StringArray
	var movieSeries pq.StringArray
	var travels pq.StringArray
	var sports pq.StringArray
	err := row.Scan(
		&newBigUser.UserId,
		&newBigUser.Alias,
		&newBigUser.Dob,
		&newBasicInfoGender,
		&newBasicInfo.FromLoc,
		&newBasicInfo.Height,
		&newBasicInfo.EducationLevel,
		&newBasicInfo.Drinking,
		&newBasicInfo.Smoking,
		&newBasicInfo.RelationshipPref,
		&newBasicInfoLookingFor,
		&newBasicInfo.Zodiac,
		&newBasicInfo.Kids,
		&newBasicInfo.Work,
		&newMatchBioId,
		&newMatchBio,
		&hobbies,
		&movieSeries,
		&travels,
		&sports,
	)
	if err != nil {
		return domain.BigUser{}, err
	}
	if newBasicInfoGender.Valid {
		newBigUser.Gender = &newBasicInfoGender.String
	}

	if newBasicInfo.FromLoc.Valid {
		newBigUser.FromLoc = &newBasicInfo.FromLoc.String
	}
	if newBasicInfo.Height.Valid {
		basicInfoHeight := int(newBasicInfo.Height.Int16)
		newBigUser.Height = &basicInfoHeight
	}
	if newBasicInfo.EducationLevel.Valid {
		newBigUser.EducationLevel = &newBasicInfo.EducationLevel.String
	}
	if newBasicInfo.Drinking.Valid {
		newBigUser.Drinking = &newBasicInfo.Drinking.String
	}
	if newBasicInfo.Smoking.Valid {
		newBigUser.Smoking = &newBasicInfo.Smoking.String
	}
	if newBasicInfo.RelationshipPref.Valid {
		newBigUser.RelationshipPref = &newBasicInfo.RelationshipPref.String
	}
	if newBasicInfoLookingFor.Valid {
		newBigUser.LookingFor = &newBasicInfoLookingFor.String
	}
	if newBasicInfo.Zodiac.Valid {
		newBigUser.Zodiac = &newBasicInfo.Zodiac.String
	}
	if newBasicInfo.Kids.Valid {
		basicInfoKids := int(newBasicInfo.Kids.Int16)
		newBigUser.Kids = &basicInfoKids
	}
	if newBasicInfo.Work.Valid {
		newBigUser.Work = &newBasicInfo.Work.String
	}
	if newMatchBioId.Valid {
		newBigUser.BioId = &newMatchBioId.String
	}
	if newMatchBio.Valid {
		newBigUser.Bio = &newMatchBio.String
	}
	for i := range hobbies {
		newBigUser.Hobbies = append(newBigUser.Hobbies, domain.Hobbie{
			Hobbie: hobbies[i],
		})
	}
	for i := range movieSeries {
		newBigUser.MovieSeries = append(newBigUser.MovieSeries, domain.MovieSerie{
			MovieSerie: movieSeries[i],
		})
	}
	for i := range travels {
		newBigUser.Travels = append(newBigUser.Travels, domain.Travel{
			Travel: travels[i],
		})
	}
	for i := range sports {
		newBigUser.Sports = append(newBigUser.Sports, domain.Sport{
			Sport: sports[i],
		})
	}
	return newBigUser, nil

}
