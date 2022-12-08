package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/util"
)

type Interest interface {
	InsertNewStats(interestId string) error

	GetInterest(userId string) (domain.Interest, error)

	InsertInterestBio(intr *domain.Bio) error
	SelectInterestBio(userId string) (domain.Bio, error)
	UpdateInterestBio(intr domain.Bio) error

	InsertInterestHobbies(interestId string, hobbies []domain.Hobbie) error
	UpdateInterestHobbies(interestId string, hobbies []domain.Hobbie) error
	DeleteInterestHobbies(interestId string, ids []string) error

	InsertInterestMovieSeries(interestId string, movieSeries []domain.MovieSerie) error
	UpdateInterestMovieSeries(interestId string, movieSeries []domain.MovieSerie) error
	DeleteInterestMovieSeries(interestId string, ids []string) error

	InsertInterestTraveling(interestId string, travels []domain.Travel) error
	UpdateInterestTraveling(interestId string, travels []domain.Travel) error
	DeleteInterestTraveling(interestId string, ids []string) error

	InsertInterestSports(interestId string, sports []domain.Sport) error
	UpdateInterestSport(interestId string, sports []domain.Sport) error
	DeleteInterestSports(interestId string, ids []string) error
}

func NewInterest(db *sqlx.DB) *IntrConn {
	return &IntrConn{
		conn: db,
	}
}

type IntrConn struct {
	conn *sqlx.DB
}

func (i *IntrConn) GetInterest(userId string) (domain.Interest, error) {
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
	var intr domain.Interest
	err := i.conn.GetContext(ctx, &intr.Bio, query, userId)
	if err != nil {
		err = i.parsingError(err)
		return domain.Interest{}, err
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

func (i *IntrConn) InsertInterestBio(intr *domain.Bio) error {
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
func (i *IntrConn) SelectInterestBio(userId string) (domain.Bio, error) {
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

	var bio domain.Bio
	err := i.conn.GetContext(ctx, &bio, query, userId)
	if err != nil {
		err = i.parsingError(err)
		return domain.Bio{}, err
	}
	return bio, nil
}
func (i *IntrConn) UpdateInterestBio(intr domain.Bio) error {
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
		WHERE interest_id = $2`
	statsHobbiesMinusQ = `
	UPDATE interest_statistics SET
		hobbie_count = hobbie_count - $1
		WHERE interest_id = $2`
)

func (i *IntrConn) InsertInterestHobbies(interestId string, hobbies []domain.Hobbie) error {
	stmt := ``
	args := make([]any, 0, len(hobbies))
	args = append(args, interestId)

	for i, val := range hobbies {
		stmt += fmt.Sprintf("($%d, $%d),", 1, i+2)
		args = append(args, val.Hobbie)
	}
	stmt = stmt[:len(stmt)-1]
	query := fmt.Sprintf(`
	INSERT INTO hobbies 
		(interest_id, hobbie)
	VALUES %s RETURNING id`, stmt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTx(ctx, func(q *sqlx.DB) error {
		rows, err := q.QueryxContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer func(rows *sqlx.Rows) {
			err := rows.Close()
			if err != nil {
				log.Panic(err)
			}
		}(rows)

		iter := 0
		for rows.Next() {
			if err := rows.Scan(&hobbies[iter].Id); err != nil {
				return err
			}
			if err != nil {
				return err
			}

			iter++
		}
		if err := rows.Err(); err != nil {
			return err
		}
		_, err = q.ExecContext(ctx, statsHobbiesPlusQ, len(hobbies), interestId)
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
func (i *IntrConn) UpdateInterestHobbies(interestId string, hobbies []domain.Hobbie) error {
	args := make([]any, 0, len(hobbies))
	args = append(args, interestId)
	stmnt := ``
	var newHobbies int
	for i, val := range hobbies {
		p1 := i * 2
		stmnt += fmt.Sprintf("(uuid($%d::TEXT), uuid($%d::TEXT), $%d),", 1, p1+2, p1+3)
		if val.Id == "" {
			newHobbies++
			val.Id = util.RandomUUID()
		}
		args = append(args, val.Id, val.Hobbie)
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

func (i *IntrConn) DeleteInterestHobbies(interestId string, ids []string) error {
	query := `
	DELETE FROM hobbies
	WHERE id IN (`
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("$%d,", i+1)
		args = append(args, id)
	}
	query = query[:len(query)-1] + ")"

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
			return common.WrapWithNewError(fmt.Errorf("no rows affected by delete"), http.StatusUnprocessableEntity, "invalid Id in one of hobbies")
		}
		_, err = q.ExecContext(ctx, statsHobbiesMinusQ, len(ids), interestId)
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

const (
	statsMovieSeriesPlusQ = `
	UPDATE interest_statistics SET
		movie_serie_count = movie_serie_count + $1
		WHERE interest_id = $2`
	statsMovieSeriesMinusQ = `
	UPDATE interest_statistics SET
		movie_serie_count = movie_serie_count - $1
		WHERE interest_id = $2`
)

func (i *IntrConn) InsertInterestMovieSeries(interestId string, movieSeries []domain.MovieSerie) error {
	stmt := ``
	args := make([]any, 0, len(movieSeries))
	args = append(args, interestId)
	for i, val := range movieSeries {
		stmt += fmt.Sprintf("($%d, $%d),", 1, i+2)
		args = append(args, val.MovieSerie)
	}
	stmt = stmt[:len(stmt)-1]

	query := fmt.Sprintf(`
	INSERT INTO movie_series 
		(interest_id, movie_serie)
	VALUES %s RETURNING id`, stmt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTx(ctx, func(q *sqlx.DB) error {
		rows, err := q.QueryxContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer func(rows *sqlx.Rows) {
			err := rows.Close()
			if err != nil {
				log.Panic(err)
			}
		}(rows)

		iter := 0
		for rows.Next() {
			if err := rows.Scan(&movieSeries[iter].Id); err != nil {
				return err
			}
			if err != nil {
				return err
			}

			iter++
		}
		if err := rows.Err(); err != nil {
			return err
		}
		_, err = q.ExecContext(ctx, statsMovieSeriesPlusQ, len(movieSeries)-1, interestId)
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

func (i *IntrConn) UpdateInterestMovieSeries(interestId string, movieSeries []domain.MovieSerie) error {
	args := make([]any, 0, len(movieSeries))
	args = append(args, interestId)
	stmnt := ``
	var newMovies int
	for i, val := range movieSeries {
		p1 := i * 2
		stmnt += fmt.Sprintf("(uuid($%d::TEXT), uuid($%d::TEXT), $%d),", 1, p1+2, p1+3)
		if val.Id == "" {
			newMovies++
			val.Id = util.RandomUUID()
		}
		args = append(args, val.Id, val.MovieSerie)
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
func (i *IntrConn) DeleteInterestMovieSeries(interestId string, ids []string) error {
	query := `
	DELETE FROM movie_series
	WHERE id IN (`
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("$%d,", i+1)
		args = append(args, id)
	}
	query = query[:len(query)-1] + ")"

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
		_, err = q.ExecContext(ctx, statsMovieSeriesMinusQ, len(ids)-1, interestId)
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

const (
	statsTravelingPlusQ = `
	UPDATE interest_statistics SET
		traveling_count = traveling_count + $1
		WHERE interest_id = $2`
	statsTravelingMinusQ = `
	UPDATE interest_statistics SET
		traveling_count = traveling_count - $1
		WHERE interest_id = $2`
)

func (i *IntrConn) InsertInterestTraveling(interestId string, travels []domain.Travel) error {
	stmt := ``
	args := make([]any, 0, len(travels))
	args = append(args, interestId)
	for i, val := range travels {
		stmt += fmt.Sprintf("($%d, $%d),", 1, i+2)
		args = append(args, val.Travel)
	}
	stmt = stmt[:len(stmt)-1]
	query := fmt.Sprintf(`
	INSERT INTO traveling 
		(interest_id, travel)
	VALUES %s RETURNING id`, stmt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		rows, err := q.QueryxContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer func(rows *sqlx.Rows) {
			err := rows.Close()
			if err != nil {
				log.Panic(err)
			}
		}(rows)

		iter := 0
		for rows.Next() {
			if err := rows.Scan(&travels[iter].Id); err != nil {
				return err
			}
			if err != nil {
				return err
			}

			iter++
		}
		if err := rows.Err(); err != nil {
			return err
		}
		_, err = q.ExecContext(ctx, statsTravelingPlusQ, len(travels), interestId)
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

func (i *IntrConn) UpdateInterestTraveling(interestId string, travels []domain.Travel) error {
	args := make([]any, 0, len(travels))
	args = append(args, interestId)
	stmnt := ``
	var newTravel int
	for i, val := range travels {
		p1 := i * 2
		stmnt += fmt.Sprintf("(uuid($%d::TEXT), uuid($%d::TEXT), $%d),", 1, p1+2, p1+3)
		if val.Id == "" {
			newTravel++
			val.Id = util.RandomUUID()
		}
		args = append(args, val.Id, val.Travel)
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
func (i *IntrConn) DeleteInterestTraveling(interestId string, ids []string) error {
	query := `
	DELETE FROM traveling
	WHERE id IN (`
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("$%d,", i+1)
		args = append(args, id)
	}
	query = query[:len(query)-1] + ")"

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
		_, err = q.ExecContext(ctx, statsTravelingMinusQ, len(ids)-1, interestId)
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

const (
	statsSportsPlusQ = `
	UPDATE interest_statistics SET
		sport_count = sport_count + $1
		WHERE interest_id = $2`
	statsSportsMinusQ = `
	UPDATE interest_statistics SET
		sport_count = sport_count - $1
		WHERE interest_id = $2`
)

func (i *IntrConn) InsertInterestSports(interestId string, sports []domain.Sport) error {
	stmt := ``
	args := make([]any, 0, len(sports))
	args = append(args, interestId)
	for i, val := range sports {
		stmt += fmt.Sprintf("($%d, $%d),", 1, i+2)
		args = append(args, val.Sport)
	}
	stmt = stmt[:len(stmt)-1]

	query := fmt.Sprintf(`
	INSERT INTO sports 
		(interest_id, sport)
	VALUES %s RETURNING id`, stmt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTx(ctx, func(q *sqlx.DB) error {
		rows, err := q.QueryxContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer func(rows *sqlx.Rows) {
			err := rows.Close()
			if err != nil {
				log.Panic(err)
			}
		}(rows)

		iter := 0
		for rows.Next() {
			if err := rows.Scan(&sports[iter].Id); err != nil {
				return err
			}
			if err != nil {
				return err
			}

			iter++
		}
		if err := rows.Err(); err != nil {
			return err
		}
		_, err = q.ExecContext(ctx, statsSportsPlusQ, len(sports), interestId)
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
func (i *IntrConn) UpdateInterestSport(interestId string, sports []domain.Sport) error {
	args := make([]any, 0, len(sports))
	args = append(args, interestId)
	stmnt := ``
	var newSport int
	for i, val := range sports {
		p1 := i * 2
		stmnt += fmt.Sprintf("(uuid($%d::TEXT), uuid($%d::TEXT), $%d),", 1, p1+2, p1+3)
		if val.Id == "" {
			newSport++
			val.Id = util.RandomUUID()
		}
		args = append(args, val.Id, val.Sport)
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

func (i *IntrConn) DeleteInterestSports(interestId string, ids []string) error {
	query := `
	DELETE FROM sports
	WHERE id IN (`
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("$%d,", i+1)
		args = append(args, id)
	}
	query = query[:len(query)-1] + ")"

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
		_, err = q.ExecContext(ctx, statsSportsMinusQ, len(ids)-1, interestId)
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
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "unique constraint on interest")
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
