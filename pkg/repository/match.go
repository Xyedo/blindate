package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/domain/entity"
)

type Match interface {
	InsertNewMatch(fromUserId, toUserId string, reqStatus domain.MatchStatus) (string, error)
	SelectMatchReqToUserId(userId string) ([]domain.MatchUser, error)
	UpdateMatchById(matchEntity entity.Match) error
	GetMatchById(matchId string) (entity.Match, error)
}

func NewMatch(conn *sqlx.DB) *MatchConn {
	return &MatchConn{
		conn: conn,
	}
}

type MatchConn struct {
	conn *sqlx.DB
}

func (m *MatchConn) InsertNewMatch(fromUserId, toUserId string, reqStatus domain.MatchStatus) (string, error) {
	query := `
	INSERT INTO match(
		request_from, 
		request_to, 
		request_status,
		created_at
		)
	VALUES($1,$2,$3,$4)
	RETURNING id`
	args := []any{fromUserId, toUserId, string(reqStatus), time.Now()}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var matchId string
	err := m.conn.GetContext(ctx, &matchId, query, args...)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.Is(err, context.Canceled):
			return "", common.WrapError(err, common.ErrTooLongAccessingDB)
		case errors.As(err, &pqErr):
			switch pqErr.Code {
			case "23503":
				switch {
				case strings.Contains(pqErr.Constraint, "request_from"):
					return "", common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "invalid user on requestFrom")
				case strings.Contains(pqErr.Constraint, "request_to"):
					return "", common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "invalid user on requestTo")
				case strings.Contains(pqErr.Constraint, "request_status"):
					return "", common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "invalid enums on requestStatus")
				case strings.Contains(pqErr.Constraint, "reveal_status"):
					return "", common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "invalid enums on revealStatus")
				}
			case "23505":
				return "", common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "match already created")
			default:
				return "", pqErr
			}
		default:
			return "", err
		}
	}
	return matchId, nil
}

func (m *MatchConn) SelectMatchReqToUserId(userId string) ([]domain.MatchUser, error) {
	query := `
	SELECT 
		m.id as match_id,
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
	FROM match m
	RIGHT JOIN users u
		ON u.id = m.request_from
	LEFT JOIN basic_info b
		ON u.id = b.user_id
	LEFT JOIN interests i
		ON i.user_id = u.id
	WHERE m.request_to = $1
		AND m.request_status = 'requested'
	ORDER BY m.created_at ASC
	LIMIT 20`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.conn.QueryxContext(ctx, query, userId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, common.WrapError(err, common.ErrTooLongAccessingDB)
		}

		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}(rows)
	matchs := make([]domain.MatchUser, 0)
	for rows.Next() {
		newMatch, err := m.createCandidatematch(rows)
		if err != nil {
			return nil, err
		}
		matchs = append(matchs, newMatch)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return matchs, nil

}

func (m *MatchConn) GetMatchById(matchId string) (entity.Match, error) {
	query := `
		SELECT
			id,
			request_from,
			request_to,
			request_status,
			created_at,
			accepted_at,
			reveal_status,
			revealed_at
		FROM match
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var matchEntity entity.Match
	err := m.conn.GetContext(ctx, &matchEntity, query, matchId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return entity.Match{}, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Match{}, common.WrapError(err, common.ErrResourceNotFound)
		}
		return entity.Match{}, err
	}
	return matchEntity, err
}
func (m *MatchConn) UpdateMatchById(matchEntity entity.Match) error {
	query := `
	UPDATE match SET
		request_status=$1, 
		accepted_at=$2, 
		reveal_status=$3, 
		revealed_at=$4
	WHERE id = $5
	RETURNING id`
	args := []any{
		matchEntity.RequestStatus,
		matchEntity.AcceptedAt,
		matchEntity.RevealStatus,
		matchEntity.RevealedAt,
		matchEntity.Id,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.conn.GetContext(ctx, &matchEntity.Id, query, args...)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return common.WrapError(err, common.ErrResourceNotFound)
		case errors.Is(err, context.Canceled):
			return common.WrapError(err, common.ErrTooLongAccessingDB)
		case errors.As(err, &pqErr):
			switch pqErr.Code {
			case "23503":
				switch {
				case strings.Contains(pqErr.Constraint, "request_status"):
					return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "invalid enums on requestStatus")
				case strings.Contains(pqErr.Constraint, "reveal_status"):
					return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "invalid enums on revealStatus")
				}
			}
		default:
			return err
		}
	}
	return nil
}
func (*MatchConn) createCandidatematch(row sqlx.ColScanner) (domain.MatchUser, error) {
	var newMatch domain.MatchUser
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
		&newMatch.MatchId,
		&newMatch.UserId,
		&newMatch.Alias,
		&newMatch.Dob,
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
		return domain.MatchUser{}, err
	}
	if newBasicInfoGender.Valid {
		newMatch.Gender = &newBasicInfoGender.String
	}

	if newBasicInfo.FromLoc.Valid {
		newMatch.FromLoc = &newBasicInfo.FromLoc.String
	}
	if newBasicInfo.Height.Valid {
		basicInfoHeight := int(newBasicInfo.Height.Int16)
		newMatch.Height = &basicInfoHeight
	}
	if newBasicInfo.EducationLevel.Valid {
		newMatch.EducationLevel = &newBasicInfo.EducationLevel.String
	}
	if newBasicInfo.Drinking.Valid {
		newMatch.Drinking = &newBasicInfo.Drinking.String
	}
	if newBasicInfo.Smoking.Valid {
		newMatch.Smoking = &newBasicInfo.Smoking.String
	}
	if newBasicInfo.RelationshipPref.Valid {
		newMatch.RelationshipPref = &newBasicInfo.RelationshipPref.String
	}
	if newBasicInfoLookingFor.Valid {
		newMatch.LookingFor = &newBasicInfoLookingFor.String
	}
	if newBasicInfo.Zodiac.Valid {
		newMatch.Zodiac = &newBasicInfo.Zodiac.String
	}
	if newBasicInfo.Kids.Valid {
		basicInfoKids := int(newBasicInfo.Kids.Int16)
		newMatch.Kids = &basicInfoKids
	}
	if newBasicInfo.Work.Valid {
		newMatch.Work = &newBasicInfo.Work.String
	}
	if newMatchBioId.Valid {
		newMatch.BioId = &newMatchBioId.String
	}
	if newMatchBio.Valid {
		newMatch.Bio = &newMatchBio.String
	}
	for i := range hobbies {
		newMatch.Hobbies = append(newMatch.Hobbies, domain.Hobbie{
			Hobbie: hobbies[i],
		})
	}
	for i := range movieSeries {
		newMatch.MovieSeries = append(newMatch.MovieSeries, domain.MovieSerie{
			MovieSerie: movieSeries[i],
		})
	}
	for i := range travels {
		newMatch.Travels = append(newMatch.Travels, domain.Travel{
			Travel: travels[i],
		})
	}
	for i := range sports {
		newMatch.Sports = append(newMatch.Sports, domain.Sport{
			Sport: sports[i],
		})
	}
	return newMatch, nil
}
