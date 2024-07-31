package types

import "time"

type CourtCase struct {
	Cnj           string    `json:"cnj"`
	Plaintiff     string    `json:"plaintiff"`
	Defendant     string    `json:"defendant"`
	CourtOfOrigin string    `json:"courtOfOrigin"`
	StartDate     time.Time `json:"startDate"`
}
