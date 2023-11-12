package pg

import "github.com/jmoiron/sqlx"

func In(query string, args ...any) (string, []any, error) {
	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return "", nil, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	return query, args, nil
}
