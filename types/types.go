package types

import "time"

type CourtCase struct {
	Cnj           string    `json:"cnj"`
	Plaintiff     string    `json:"plaintiff"`
	Defendant     string    `json:"defendant"`
	CourtOfOrigin string    `json:"court_of_origin"`
	StartDate     time.Time `json:"start_date"`
}

type CaseUpdate struct {
	UpdateDate    time.Time `json:"update_date"`
	UpdateDetails string    `json:"update_details"`
}
