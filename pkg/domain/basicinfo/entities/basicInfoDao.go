package basicInfoEntity

import (
	"database/sql"
	"time"
)

type DAO struct {
	UserId           string         `db:"user_id"`
	Gender           string         `db:"gender"`
	FromLoc          sql.NullString `db:"from_loc"`
	Height           sql.NullInt16  `db:"height"`
	EducationLevel   sql.NullString `db:"education_level"`
	Drinking         sql.NullString `db:"drinking"`
	Smoking          sql.NullString `db:"smoking"`
	RelationshipPref sql.NullString `db:"relationship_pref"`
	LookingFor       string         `db:"looking_for"`
	Zodiac           sql.NullString `db:"zodiac"`
	Kids             sql.NullInt16  `db:"kids"`
	Work             sql.NullString `db:"work"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}
