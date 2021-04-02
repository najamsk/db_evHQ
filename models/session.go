package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	Base
	Title   string
	Summary string
	Details string
	// Poster           string
	// Thumbnail        string
	IsActive         bool
	StartDate        time.Time
	StartDateDisplay string
	EndDate          time.Time
	EndDateDisplay   string
	DurationDisplay  string
	Address          string
	Venue            string
	Seats            int
	IsFeatured		 bool
	SortOrder		 int
	
	GeoLocation      GeoLocation   `gorm:"embedded"`
	Speakers         []*Speaker    `gorm:"many2many:session_speakers;"`
	Ambassadors      []*Ambassador `gorm:"many2many:session_ambassadors;"`
	SessionTickets   []SessionTicket
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID


	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}
