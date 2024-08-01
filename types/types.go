package types

import (
	"errors"
	"time"
)

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

type ErrResponse struct {
	Error string `json:"error"`
}

func (c *CourtCase) Validate() error {
	if c.Cnj == "" {
		return errors.New("cnj field cannot be empty")
	}
	if c.Plaintiff == "" {
		return errors.New("plaintiff field cannot be empty")
	}
	if c.Defendant == "" {
		return errors.New("defendant field cannot be empty")
	}
	if c.CourtOfOrigin == "" {
		return errors.New("court_of_origin field cannot be empty")
	}
	if c.StartDate.IsZero() {
		return errors.New("start_date field must be a valid date")
	}
	return nil


}
