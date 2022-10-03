package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/entity"
)

func NewInterest(db *sqlx.DB) *interest {
	return &interest{
		db,
	}
}

type interest struct {
	*sqlx.DB
}

func (i *interest) InsertInterest(intr *entity.Interest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := execTx(i.DB, ctx, func(q *sqlx.DB) error {
		q1 := `
		INSERT INTO interests (
			user_id, 
			bio, 
			created_at, 
			updated_at)
		VALUES($1,$2,$3,$3) 
		RETURNING id`
		err := q.SelectContext(ctx, &intr.Id, q1, intr.UserId, intr.Bio, time.Now())
		if err != nil {
			return err
		}
		if len(intr.Hobbies) > 0 {
			err = bulkInsertHobbies(q, ctx, intr.Id, intr.Hobbies)
			if err != nil {
				return err
			}
		}
		if len(intr.MoviesSeries) > 0 {
			err = bulkInsertMovieSeries(q, ctx, intr.Id, intr.MoviesSeries)
			if err != nil {
				return err
			}
		}
		if len(intr.Traveling) > 0 {
			err = bulkInsertTraveling(q, ctx, intr.Id, intr.Traveling)
			if err != nil {
				return err
			}
		}
		if len(intr.Sports) > 0 {
			err = bulkInsertSport(q, ctx, intr.Id, intr.Sports)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (i *interest) GetInterest(userId string) (*entity.Interest, error) {
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

	var intr entity.Interest
	err := i.GetContext(ctx, &intr, query, userId)
	if err != nil {
		return nil, err
	}

	query = `SELECT id, hobbie FROM hobbies WHERE interest_id = $1`
	err = i.SelectContext(ctx, &intr.Hobbies, query, intr.Id)
	if err != nil {
		return &intr, err
	}

	query = `SELECT id, movie_serie FROM movie_series WHERE interest_id = $1`
	err = i.SelectContext(ctx, &intr.MoviesSeries, query, intr.Id)
	if err != nil {
		return &intr, err
	}

	query = `SELECT id, travel FROM traveling WHERE interest_id = $1`
	err = i.SelectContext(ctx, &intr.Traveling, query, intr.Id)
	if err != nil {
		return &intr, err
	}

	query = `SELECT sport FROM sports WHERE interest_id = $1`
	err = i.SelectContext(ctx, &intr.Sports, query, intr.Id)
	if err != nil {
		return &intr, err
	}
	return &intr, nil
}
func (i *interest) UpdateInterest(intr *entity.Interest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := execTx(i.DB, ctx, func(q *sqlx.DB) error {
		query := `UPDATE interests SET bio = $1, updated_at=$2  WHERE user_id = $3 RETURNING id`
		err := q.GetContext(ctx, &intr.Id, query, intr.Bio, time.Now(), intr.UserId)
		if err != nil {
			return err
		}
		if len(intr.Hobbies) > 0 {
			stmnt := ``
			for _, val := range intr.Hobbies {
				stmnt += fmt.Sprintf("(%s, %s),", val.Id, val.Hobbie)
			}
			stmnt = stmnt[:len(stmnt)-1]

			query = `
			UPDATE hobbies as h SET
				hobbie = h2.hobbie
			FROM (VALUES
				$1	
			) as h2(id,hobbie)
			WHERE h2.id = h.id`

			_, err = q.ExecContext(ctx, query, stmnt)
			if err != nil {
				return err
			}
		}
		if len(intr.MoviesSeries) > 0 {
			stmnt := ``
			for _, val := range intr.MoviesSeries {
				stmnt += fmt.Sprintf("(%s, %s),", val.Id, val.MovieSerie)
			}
			stmnt = stmnt[:len(stmnt)-1]
			query = `
			UPDATE movie_series as m SET
				movie_serie = m2.movie_serie
			FROM (VALUES
				$1
			) AS m2(id, movie_serie)
			WHERE m2.id = m.id`
			_, err = q.ExecContext(ctx, query, stmnt)
			if err != nil {
				return err
			}
		}
		if len(intr.Traveling) > 0 {
			stmnt := ``
			for _, val := range intr.Traveling {
				stmnt += fmt.Sprintf("(%s, %s),", val.Id, val.Travel)
			}
			stmnt = stmnt[:len(stmnt)-1]

			query = `
			UPDATE traveling as t SET
				travel = t2.travel
			FROM (VALUES
				$1
			) AS t2(id, travel)
			WHERE t2.id = t.id`
			_, err = q.ExecContext(ctx, query, stmnt)
			if err != nil {
				return err
			}
		}
		if len(intr.Sports) > 0 {
			stmnt := ``
			for _, val := range intr.Sports {
				stmnt += fmt.Sprintf("(%s, %s),", val.Id, val.Sport)
			}
			stmnt = stmnt[:len(stmnt)-1]

			query = `
			UPDATE sports as s SET
				sport = s2.sport
			FROM (VALUES
				$1
			) AS s2(id, sport)
			WHERE s2.id = s.id`
			_, err = q.ExecContext(ctx, query, stmnt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}
func bulkInsertSport(q *sqlx.DB, ctx context.Context, interestId string, sports []entity.Sports) error {
	query := `
	INSERT INTO sports 
		(interest_id, sport)
	VALUES `

	args := make([]any, 0, len(sports))
	args = append(args, interestId)
	for i, val := range sports {
		query += fmt.Sprintf("($%d, $%d),", 1, i+2)
		args = append(args, val.Sport)
	}
	query = query[:len(query)-1]
	_, err := q.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func bulkInsertTraveling(q *sqlx.DB, ctx context.Context, interestId string, traveling []entity.Traveling) error {
	query := `
	INSERT INTO traveling 
		(interest_id, travel)
	VALUES `

	args := make([]any, 0, len(traveling))
	args = append(args, interestId)
	for i, val := range traveling {
		query += fmt.Sprintf("($%d, $%d),", 1, i+2)
		args = append(args, val.Travel)
	}
	query = query[:len(query)-1]
	_, err := q.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func bulkInsertMovieSeries(q *sqlx.DB, ctx context.Context, interestId string, moviesSeries []entity.MovieSeries) error {
	query := `
	INSERT INTO movie_series 
		(interest_id, movie_serie)
	VALUES `

	args := make([]any, 0, len(moviesSeries))
	args = append(args, interestId)
	for i, val := range moviesSeries {
		query += fmt.Sprintf("($%d, $%d),", 1, i+2)
		args = append(args, val.MovieSerie)
	}
	query = query[:len(query)-1]
	_, err := q.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func bulkInsertHobbies(q *sqlx.DB, ctx context.Context, interestId string, hobies []entity.Hobies) error {
	query := `
		INSERT INTO hobbies 
			(interest_id, hobbie)
		VALUES `

	args := make([]any, 0, len(hobies))
	args = append(args, interestId)
	for i, val := range hobies {
		query += fmt.Sprintf("($%d, $%d),", 1, i+2)
		args = append(args, val.Hobbie)
	}
	query = query[:len(query)-1]
	_, err := q.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
