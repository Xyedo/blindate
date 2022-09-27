package repository

import "github.com/jmoiron/sqlx"

func NewBasicInfo(db *sqlx.DB) *basicInfo {
	return &basicInfo{
		db,
	}
}

type basicInfo struct {
	*sqlx.DB
}
