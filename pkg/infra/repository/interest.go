package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/common"
	interestEntity "github.com/xyedo/blindate/pkg/domain/interest/entities"
	"github.com/xyedo/blindate/pkg/util"
)

func NewInterest(db *sqlx.DB) *IntrConn {
	return &IntrConn{
		conn: db,
	}
}

type IntrConn struct {
	conn *sqlx.DB
}

func (i *IntrConn) GetInterest(userId string) (interestEntity.FullDTO, error) {
	query := `
	SELECT
		id, 
		user_id, 
		bio, 
		created_at, 
		updated_at 
	FROM interests
	WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var intr interestEntity.FullDTO
	err := i.conn.GetContext(ctx, &intr.Bio, query, userId)
	if err != nil {
		err = i.parsingError(err)
		return interestEntity.FullDTO{}, err
	}
	query = `SELECT id, hobbie FROM hobbies WHERE interest_id = $1`
	err = i.conn.SelectContext(ctx, &intr.Hobbies, query, intr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			err = i.parsingError(err)
			return intr, err
		}
	}

	query = `SELECT id, movie_serie FROM movie_series WHERE interest_id = $1`
	err = i.conn.SelectContext(ctx, &intr.MovieSeries, query, intr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			err = i.parsingError(err)
			return intr, err
		}
	}
	query = `SELECT id, travel FROM traveling WHERE interest_id = $1`
	err = i.conn.SelectContext(ctx, &intr.Travels, query, intr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			err = i.parsingError(err)
			return intr, err
		}

	}
	query = `SELECT id, sport FROM sports WHERE interest_id = $1`
	err = i.conn.SelectContext(ctx, &intr.Sports, query, intr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			err = i.parsingError(err)
			return intr, err
		}
	}
	return intr, nil
}

func (i *IntrConn) InsertInterestBio(intr *interestEntity.BioDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q1 := `
		INSERT INTO interests (
			user_id, 
			bio, 
			created_at, 
			updated_at)
		VALUES($1,$2,$3,$3) 
		RETURNING id`
	err := i.conn.GetContext(ctx, &intr.Id, q1, intr.UserId, intr.Bio, time.Now())
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *IntrConn) InsertNewStats(interestId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `
	INSERT INTO interest_statistics(interest_id)
	VALUES($1)`

	_, err := i.conn.ExecContext(ctx, q, interestId)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *IntrConn) SelectInterestBio(userId string) (interestEntity.BioDTO, error) {
	query := `
	SELECT
		id, 
		user_id, 
		bio, 
		created_at, 
		updated_at 
	FROM interests
	WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var bio interestEntity.BioDTO
	err := i.conn.GetContext(ctx, &bio, query, userId)
	if err != nil {
		err = i.parsingError(err)
		return interestEntity.BioDTO{}, err
	}
	return bio, nil
}
func (i *IntrConn) UpdateInterestBio(intr interestEntity.BioDTO) error {
	query := `UPDATE interests SET bio = $1, updated_at=$2  WHERE user_id = $3 RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var retId string
	err := i.conn.GetContext(ctx, &retId, query, intr.Bio, time.Now(), intr.UserId)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}

const (
	statsHobbiesPlusQ = `
	UPDATE interest_statistics SET
		hobbie_count = hobbie_count + $1
	WHERE interest_id = $2
	RETURNING interest_id`
	statsHobbiesMinusQ = `
	UPDATE interest_statistics SET
		hobbie_count = hobbie_count - $1
	WHERE interest_id = $2
	RETURNING interest_id`
)

func (i *IntrConn) InsertInterestHobbies(interestId string, hobbies []interestEntity.HobbieDTO) error {
	stmt := ``
	args := make([]any, 0, len(hobbies))
	args = append(args, interestId)

	for i := range hobbies {
		p1 := i * 2
		stmt += fmt.Sprintf("($%d, $%d::uuid, $%d),", 1, p1+2, p1+3)
		newUid := util.RandomUUID()
		args = append(args, newUid, hobbies[i].Hobbie)
		hobbies[i].Id = newUid
	}
	stmt = stmt[:len(stmt)-1]
	query := fmt.Sprintf(`
	INSERT INTO hobbies 
		(interest_id, id, hobbie)
	VALUES %s RETURNING id`, stmt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTx(ctx, func(q *sqlx.DB) error {
		var retIds []string
		err := q.SelectContext(ctx, &retIds, query, args...)
		if err != nil {
			return err
		}
		var retInterestId string
		err = q.GetContext(ctx, &retInterestId, statsHobbiesPlusQ, len(retIds), interestId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *IntrConn) UpdateInterestHobbies(interestId string, hobbies []interestEntity.HobbieDTO) error {
	args := make([]any, 0, len(hobbies))
	args = append(args, interestId)
	stmnt := ``
	var newHobbies int
	for i := range hobbies {
		p1 := i * 2
		stmnt += fmt.Sprintf("(uuid($%d::TEXT), uuid($%d::TEXT), $%d),", 1, p1+2, p1+3)
		if hobbies[i].Id == "" {
			newHobbies++
			hobbies[i].Id = util.RandomUUID()
		}
		args = append(args, hobbies[i].Id, hobbies[i].Hobbie)
	}
	stmnt = stmnt[:len(stmnt)-1]

	query :=
		fmt.Sprintf(`
		WITH new_values(interest_id, id,hobbie) AS (
			VALUES
				%s
		),
		UPSERT AS
		(
			UPDATE hobbies h
				SET hobbie = nv.hobbie
			FROM new_values nv
			WHERE h.id =nv.id
			RETURNING h.*
		)
		INSERT INTO hobbies(interest_id,hobbie)
		SELECT interest_id, hobbie 
		FROM new_values
		WHERE NOT EXISTS (
			SELECT 1
			FROM upsert up
			WHERE up.id = new_values.id)`, stmnt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		res, err := q.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if ro == int64(0) {
			return common.WrapWithNewError(fmt.Errorf("no rows affected by delete"), http.StatusUnprocessableEntity, "invalid Id in one of hobbies")
		}

		_, err = q.ExecContext(ctx, statsHobbiesPlusQ, newHobbies, interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return i.parsingError(err)
	}

	return nil
}

func (i *IntrConn) DeleteInterestHobbies(interestId string, ids []string) ([]string, error) {
	query := `
	DELETE FROM hobbies
	WHERE id IN (`
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("$%d,", i+1)
		args = append(args, id)
	}
	query = query[:len(query)-1] + ") RETURNING id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var retIds []string
	err := i.execTx(ctx, func(q *sqlx.DB) error {

		err := q.SelectContext(ctx, &retIds, query, args...)
		if err != nil {
			return err
		}
		var retInterestId string
		err = q.GetContext(ctx, &retInterestId, statsHobbiesMinusQ, len(retIds), interestId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, i.parsingError(err)
	}
	return retIds, nil
}

const (
	statsMovieSeriesPlusQ = `
	UPDATE interest_statistics SET
		movie_serie_count = movie_serie_count + $1
	WHERE interest_id = $2
	RETURNING interest_id`
	statsMovieSeriesMinusQ = `
	UPDATE interest_statistics SET
		movie_serie_count = movie_serie_count - $1
	WHERE interest_id = $2
	RETURNING interest_id`
)

func (i *IntrConn) InsertInterestMovieSeries(interestId string, movieSeries []interestEntity.MovieSerieDTO) error {
	stmt := ``
	args := make([]any, 0, len(movieSeries))
	args = append(args, interestId)
	for i := range movieSeries {
		p1 := i * 2
		stmt += fmt.Sprintf("($%d, $%d::uuid, $%d),", 1, p1+2, p1+3)
		newUid := util.RandomUUID()
		args = append(args, newUid, movieSeries[i].MovieSerie)
		movieSeries[i].Id = newUid
	}
	stmt = stmt[:len(stmt)-1]

	query := fmt.Sprintf(`
	INSERT INTO movie_series 
		(interest_id, id, movie_serie)
	VALUES %s RETURNING id`, stmt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTx(ctx, func(q *sqlx.DB) error {
		var retIds []string
		err := q.SelectContext(ctx, &retIds, query, args...)
		if err != nil {
			return err
		}
		var retInterestId string
		err = q.GetContext(ctx, &retInterestId, statsMovieSeriesPlusQ, len(retIds), interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}

func (i *IntrConn) UpdateInterestMovieSeries(interestId string, movieSeries []interestEntity.MovieSerieDTO) error {
	args := make([]any, 0, len(movieSeries))
	args = append(args, interestId)
	stmnt := ``
	var newMovies int
	for i := range movieSeries {
		p1 := i * 2
		stmnt += fmt.Sprintf("(uuid($%d::TEXT), uuid($%d::TEXT), $%d),", 1, p1+2, p1+3)
		if movieSeries[i].Id == "" {
			newMovies++
			movieSeries[i].Id = util.RandomUUID()
		}
		args = append(args, movieSeries[i].Id, movieSeries[i].MovieSerie)
	}
	stmnt = stmnt[:len(stmnt)-1]

	query :=
		fmt.Sprintf(`
		WITH new_values(interest_id, id, movie_serie) AS (
			VALUES
				%s
		),
		UPSERT AS
		(
			UPDATE movie_series m
				SET movie_serie = nv.movie_serie
			FROM new_values nv
			WHERE m.id =nv.id
			RETURNING m.*
		)
		INSERT INTO movie_series(interest_id,movie_serie)
		SELECT interest_id, movie_serie 
		FROM new_values
		WHERE NOT EXISTS (
			SELECT 1
			FROM upsert up
			WHERE up.id = new_values.id)`, stmnt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		res, err := q.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		row, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if int(row) == 0 {
			return common.WrapWithNewError(fmt.Errorf("no rows affected by delete"), http.StatusUnprocessableEntity, "invalid Id in one of movieSeries")
		}
		_, err = q.ExecContext(ctx, statsMovieSeriesPlusQ, newMovies, interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *IntrConn) DeleteInterestMovieSeries(interestId string, ids []string) ([]string, error) {
	query := `
	DELETE FROM movie_series
	WHERE id IN (`
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("$%d,", i+1)
		args = append(args, id)
	}
	query = query[:len(query)-1] + ") RETURNING id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var retIds []string
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		err := q.SelectContext(ctx, &retIds, query, args...)
		if err != nil {
			return err
		}
		var retInterestId string
		err = q.GetContext(ctx, &retInterestId, statsMovieSeriesMinusQ, len(retIds), interestId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, i.parsingError(err)
	}
	return retIds, nil
}

const (
	statsTravelingPlusQ = `
	UPDATE interest_statistics SET
		traveling_count = traveling_count + $1
		WHERE interest_id = $2
	RETURNING interest_id`
	statsTravelingMinusQ = `
	UPDATE interest_statistics SET
		traveling_count = traveling_count - $1
		WHERE interest_id = $2
	RETURNING interest_id`
)

func (i *IntrConn) InsertInterestTraveling(interestId string, travels []interestEntity.TravelDTO) error {
	stmt := ``
	args := make([]any, 0, len(travels))
	args = append(args, interestId)
	for i := range travels {
		p1 := i * 2
		stmt += fmt.Sprintf("($%d, $%d::uuid, $%d),", 1, p1+2, p1+3)
		newUid := util.RandomUUID()
		args = append(args, newUid, travels[i].Travel)
		travels[i].Id = newUid
	}
	stmt = stmt[:len(stmt)-1]
	query := fmt.Sprintf(`
	INSERT INTO traveling 
		(interest_id, id, travel)
	VALUES %s RETURNING id`, stmt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		var retIds []string
		err := q.SelectContext(ctx, &retIds, query, args...)
		if err != nil {
			return err
		}
		var retInterestId string
		err = q.GetContext(ctx, &retInterestId, statsTravelingPlusQ, len(retIds), interestId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}

func (i *IntrConn) UpdateInterestTraveling(interestId string, travels []interestEntity.TravelDTO) error {
	args := make([]any, 0, len(travels))
	args = append(args, interestId)
	stmnt := ``
	var newTravel int
	for i := range travels {
		p1 := i * 2
		stmnt += fmt.Sprintf("(uuid($%d::TEXT), uuid($%d::TEXT), $%d),", 1, p1+2, p1+3)
		if travels[i].Id == "" {
			newTravel++
			travels[i].Id = util.RandomUUID()
		}
		args = append(args, travels[i].Id, travels[i].Travel)
	}
	stmnt = stmnt[:len(stmnt)-1]

	query :=
		fmt.Sprintf(`
		WITH new_values(interest_id, id,travel) AS (
			VALUES
				%s
		),
		UPSERT AS
		(
			UPDATE traveling t
				SET travel = nv.travel
			FROM new_values nv
			WHERE t.id =nv.id
			RETURNING t.*
		)
		INSERT INTO traveling(interest_id,travel)
		SELECT interest_id,travel FROM new_values
		WHERE NOT EXISTS (
			SELECT 1
			FROM upsert up
			WHERE up.id = new_values.id)`, stmnt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		res, err := q.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		row, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if int(row) == 0 {
			return common.WrapWithNewError(fmt.Errorf("no rows affected by delete"), http.StatusUnprocessableEntity, "invalid Id in one of travels")
		}

		_, err = q.ExecContext(ctx, statsTravelingPlusQ, newTravel, interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return i.parsingError(err)
	}

	return nil
}
func (i *IntrConn) DeleteInterestTraveling(interestId string, ids []string) ([]string, error) {
	query := `
	DELETE FROM traveling
	WHERE id IN (`
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("$%d,", i+1)
		args = append(args, id)
	}
	query = query[:len(query)-1] + ") RETURNING id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var retIds []string
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		err := q.SelectContext(ctx, &retIds, query, args...)
		if err != nil {
			return err
		}
		var retInterestId string
		err = q.GetContext(ctx, &retInterestId, statsTravelingMinusQ, len(retIds), interestId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, i.parsingError(err)
	}
	return retIds, nil
}

const (
	statsSportsPlusQ = `
	UPDATE interest_statistics SET
		sport_count = sport_count + $1
		WHERE interest_id = $2
	RETURNING interest_id`
	statsSportsMinusQ = `
	UPDATE interest_statistics SET
		sport_count = sport_count - $1
		WHERE interest_id = $2
	RETURNING interest_id`
)

func (i *IntrConn) InsertInterestSports(interestId string, sports []interestEntity.SportDTO) error {
	stmt := ``
	args := make([]any, 0, len(sports))
	args = append(args, interestId)
	for i := range sports {
		p1 := i * 2
		stmt += fmt.Sprintf("($%d, $%d::uuid, $%d),", 1, p1+2, p1+3)
		newUid := util.RandomUUID()
		args = append(args, newUid, sports[i].Sport)
		sports[i].Id = newUid
	}
	stmt = stmt[:len(stmt)-1]

	query := fmt.Sprintf(`
	INSERT INTO sports 
		(interest_id, id, sport)
	VALUES %s RETURNING id`, stmt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTx(ctx, func(q *sqlx.DB) error {
		var retIds []string
		err := q.SelectContext(ctx, &retIds, query, args...)
		if err != nil {
			return err
		}
		var retInterestId string
		err = q.GetContext(ctx, &retInterestId, statsSportsPlusQ, len(retIds), interestId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *IntrConn) UpdateInterestSport(interestId string, sports []interestEntity.SportDTO) error {
	args := make([]any, 0, len(sports))
	args = append(args, interestId)
	stmnt := ``
	var newSport int
	for i := range sports {
		p1 := i * 2
		stmnt += fmt.Sprintf("(uuid($%d::TEXT), uuid($%d::TEXT), $%d),", 1, p1+2, p1+3)
		if sports[i].Id == "" {
			newSport++
			sports[i].Id = util.RandomUUID()
		}
		args = append(args, sports[i].Id, sports[i].Sport)
	}
	stmnt = stmnt[:len(stmnt)-1]

	query :=
		fmt.Sprintf(`
		WITH new_values(interest_id, id,sport) AS (
			VALUES
				%s
		),
		UPSERT AS
		(
			UPDATE sports s
				SET sport = nv.sport
			FROM new_values nv
			WHERE s.id =nv.id
			RETURNING s.*
		)
		INSERT INTO sports(interest_id,sport)
		SELECT interest_id,sport FROM new_values
		WHERE NOT EXISTS (
			SELECT 1
			FROM upsert up
			WHERE up.id = new_values.id)`, stmnt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTx(ctx, func(q *sqlx.DB) error {
		res, err := q.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		row, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if int(row) == 0 {
			return common.WrapWithNewError(fmt.Errorf("no rows affected by delete"), http.StatusUnprocessableEntity, "invalid Id in one of sports")
		}

		_, err = q.ExecContext(ctx, statsSportsPlusQ, newSport, interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return i.parsingError(err)
	}

	return nil
}

func (i *IntrConn) DeleteInterestSports(interestId string, ids []string) ([]string, error) {
	query := `
	DELETE FROM sports
	WHERE id IN (`
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("$%d,", i+1)
		args = append(args, id)
	}
	query = query[:len(query)-1] + ") RETURNING id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var retIds []string
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		err := q.SelectContext(ctx, &retIds, query, args...)
		if err != nil {
			return err
		}
		var retInterestId string
		err = q.GetContext(ctx, &retInterestId, statsSportsMinusQ, len(retIds), interestId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, i.parsingError(err)
	}
	return retIds, nil
}
func (*IntrConn) parsingError(err error) error {
	var pqErr *pq.Error
	switch {
	case errors.Is(err, context.Canceled):
		return common.WrapError(err, common.ErrTooLongAccessingDB)
	case errors.Is(err, sql.ErrNoRows):
		return common.WrapError(err, common.ErrResourceNotFound)
	case errors.As(err, &pqErr):
		switch pqErr.Code {
		case "23503":
			if strings.Contains(pqErr.Constraint, "interest_id") {
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "interestId is invalid")
			}
			if strings.Contains(pqErr.Constraint, "user_id") {
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "userId is not invalid")
			}
		case "23505":
			if strings.Contains(pqErr.Constraint, "interests_user_id_key") {
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "interest with this user id is already created")
			}
			if strings.Contains(pqErr.Constraint, "hobbie_unique") {
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "every hobbies must be unique")
			}
			if strings.Contains(pqErr.Constraint, "movie_serie_unique") {
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "every moviesSeries must be unique")
			}
			if strings.Contains(pqErr.Constraint, "travel_unique") {
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "every travels must be unique")
			}
			if strings.Contains(pqErr.Constraint, "sport_unique") {
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "every sports must be unique")
			}
			if strings.Contains(pqErr.Constraint, "user_id") {
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "interest already created")
			}

		case "23514":
			if strings.Contains(pqErr.Constraint, "hobbie_count") {
				return common.WrapWithNewError(err, http.StatusUnprocessableEntity, "hobbies must less than 10")
			}
			if strings.Contains(pqErr.Constraint, "movie_serie_count") {
				return common.WrapWithNewError(err, http.StatusUnprocessableEntity, "movieSeries must less than 10")
			}
			if strings.Contains(pqErr.Constraint, "traveling_count") {
				return common.WrapWithNewError(err, http.StatusUnprocessableEntity, "travels must less than 10")
			}
			if strings.Contains(pqErr.Constraint, "sport_count") {
				return common.WrapWithNewError(err, http.StatusUnprocessableEntity, "sports must less than 10")
			}
		default:
			return pqErr
		}
	}
	return err
}
func (i *IntrConn) execTx(ctx context.Context, q func(q *sqlx.DB) error) error {
	return execGeneric(i.conn, ctx, q, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}
