package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/util"
)

type Interest interface {
	InsertNewStats(interestId string) error

	GetInterest(userId string) (*domain.Interest, error)

	InsertInterestBio(intr *domain.Bio) error
	SelectInterestBio(userId string) (*domain.Bio, error)
	UpdateInterestBio(intr *domain.Bio) error

	InsertInterestHobbies(interestId string, hobbies []domain.Hobbie) error
	UpdateInterestHobbies(interestId string, hobbies []domain.Hobbie) (int64, error)
	DeleteInterestHobbies(interestId string, ids []string) (int64, error)

	InsertInterestMovieSeries(interestId string, movieSeries []domain.MovieSerie) error
	UpdateInterestMovieSeries(interestId string, movieSeries []domain.MovieSerie) (int64, error)
	DeleteInterestMovieSeries(interestId string, ids []string) (int64, error)

	InsertInterestTraveling(interestId string, travels []domain.Travel) error
	UpdateInterestTraveling(interestId string, travels []domain.Travel) (int64, error)
	DeleteInterestTraveling(interestId string, ids []string) (int64, error)

	InsertInterestSports(interestId string, sports []domain.Sport) error
	UpdateInterestSport(interestId string, sports []domain.Sport) (int64, error)
	DeleteInterestSports(interestId string, ids []string) (int64, error)
}

func NewInterest(db *sqlx.DB) *interest {
	return &interest{
		conn: db,
	}
}

type interest struct {
	conn *sqlx.DB
}

func (i *interest) GetInterest(userId string) (*domain.Interest, error) {
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
		return nil, err
	}
	query = `SELECT id, hobbie FROM hobbies WHERE interest_id = $1`
	err = i.conn.SelectContext(ctx, &intr.Hobbies, query, intr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			return &intr, err
		}
	}

	query = `SELECT id, movie_serie FROM movie_series WHERE interest_id = $1`
	err = i.conn.SelectContext(ctx, &intr.MovieSeries, query, intr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			return &intr, err
		}
	}
	query = `SELECT id, travel FROM traveling WHERE interest_id = $1`
	err = i.conn.SelectContext(ctx, &intr.Travels, query, intr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			return &intr, err
		}

	}
	query = `SELECT id, sport FROM sports WHERE interest_id = $1`
	err = i.conn.SelectContext(ctx, &intr.Sports, query, intr.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		} else {
			return &intr, err
		}
	}
	return &intr, nil
}

func (i *interest) InsertInterestBio(intr *domain.Bio) error {
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
		return err
	}
	return nil
}
func (i *interest) InsertNewStats(interestId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `
	INSERT INTO interest_statistics(interest_id)
	VALUES($1)`

	_, err := i.conn.ExecContext(ctx, q, interestId)
	if err != nil {
		return err
	}
	return nil
}
func (i *interest) SelectInterestBio(userId string) (*domain.Bio, error) {
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
		return nil, err
	}
	return &bio, nil
}
func (i *interest) UpdateInterestBio(intr *domain.Bio) error {
	query := `UPDATE interests SET bio = $1, updated_at=$2  WHERE user_id = $3 RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.conn.GetContext(ctx, &intr.Id, query, intr.Bio, time.Now(), intr.UserId)
	if err != nil {
		return err
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

func (i *interest) InsertInterestHobbies(interestId string, hobbies []domain.Hobbie) error {
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
		return err
	}
	return nil
}
func (i *interest) UpdateInterestHobbies(interestId string, hobbies []domain.Hobbie) (int64, error) {
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
	var row int64
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		res, err := q.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		row = ro

		_, err = q.ExecContext(ctx, statsHobbiesPlusQ, newHobbies, interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return row, nil
}

func (i *interest) DeleteInterestHobbies(interestId string, ids []string) (int64, error) {
	var row int64
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
		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		row = ro
		_, err = q.ExecContext(ctx, statsHobbiesMinusQ, len(ids), interestId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return row, nil
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

func (i *interest) InsertInterestMovieSeries(interestId string, movieSeries []domain.MovieSerie) error {
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
		return err
	}
	return nil
}

func (i *interest) UpdateInterestMovieSeries(interestId string, movieSeries []domain.MovieSerie) (int64, error) {
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
	var row int64
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		res, err := q.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		row = ro

		_, err = q.ExecContext(ctx, statsMovieSeriesPlusQ, newMovies, interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return row, nil
}
func (i *interest) DeleteInterestMovieSeries(interestId string, ids []string) (int64, error) {
	var row int64
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
		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		row = ro
		_, err = q.ExecContext(ctx, statsMovieSeriesMinusQ, len(ids)-1, interestId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return row, nil
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

func (i *interest) InsertInterestTraveling(interestId string, travels []domain.Travel) error {
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
		return err
	}
	return nil
}

func (i *interest) UpdateInterestTraveling(interestId string, travels []domain.Travel) (int64, error) {
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
	var row int64
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		res, err := q.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		row = ro

		_, err = q.ExecContext(ctx, statsTravelingPlusQ, newTravel, interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return row, nil
}
func (i *interest) DeleteInterestTraveling(interestId string, ids []string) (int64, error) {
	var row int64
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
		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		row = ro
		_, err = q.ExecContext(ctx, statsTravelingMinusQ, len(ids)-1, interestId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return row, nil
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

func (i *interest) InsertInterestSports(interestId string, sports []domain.Sport) error {
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
		return err
	}
	return nil
}
func (i *interest) UpdateInterestSport(interestId string, sports []domain.Sport) (int64, error) {
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
	var row int64
	err := i.execTx(ctx, func(q *sqlx.DB) error {
		res, err := q.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		row = ro

		_, err = q.ExecContext(ctx, statsSportsPlusQ, newSport, interestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return row, nil
}

func (i *interest) DeleteInterestSports(interestId string, ids []string) (int64, error) {
	var row int64
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
		ro, err := res.RowsAffected()
		if err != nil {
			return err
		}
		row = ro
		_, err = q.ExecContext(ctx, statsSportsMinusQ, len(ids)-1, interestId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return row, nil
}

func (i *interest) execTx(ctx context.Context, q func(q *sqlx.DB) error) error {
	return execGeneric(i.conn, ctx, q, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}
