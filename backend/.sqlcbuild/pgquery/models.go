// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package pgquery

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CaseUpdate struct {
	ID            pgtype.UUID
	Cnj           string
	UpdateDate    pgtype.Timestamptz
	UpdateDetails string
	CreatedAt     pgtype.Timestamptz
}

type CourtCase struct {
	ID            pgtype.UUID
	Cnj           string
	Plaintiff     string
	Defendant     string
	CourtOfOrigin string
	StartDate     pgtype.Date
	CreatedAt     pgtype.Timestamptz
}
