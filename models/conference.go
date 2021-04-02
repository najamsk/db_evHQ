package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Conference struct {
	Base
	Title   string
	Summary string
	Details string
	// Poster               string
	// Thumbnail            string
	IsActive             bool
	StartDate            time.Time
	StartDateDisplay     string
	EndDate              time.Time
	EndDateDisplay       string
	DurationDisplay      string
	Address              string
	GeoLocation          GeoLocation `gorm:"embedded"`
	ClientID             uuid.UUID
	Sessions             []Session
	FloorPlanPoster      string
	ExibitionStartupPlan string
	ExibitionSponsorPlan string

	// CreatedAt            time.Time
	// UpdatedAt            time.Time
	// DeletedAt            time.Time
}
