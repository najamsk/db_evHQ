package viewmodels

import (
	"time"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

//ConferenceEditVMRead for edit form
type ConferenceEditVMRead struct {
	Base
	ClientID        uuid.UUID
	ClientName      string
	StartDate       string
	EndDate         string
	DurationDisplay string
	Title           string
	Summary         string
	Details         string
	IsActive        bool
	Latitude        float64
	Longitude       float64
	LocationRadius  float64
	Address         string
	PosterURL		string
	ThumbnailURL 	string

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}

//ConferenceEditVMWrite maybe use this vm from controller to service insteal of too many func params.
type ConferenceEditVMWrite struct {
	ID string
	Base
	ClientID        uuid.UUID
	ClientName      string
	StartDate       time.Time
	EndDate         time.Time
	DurationDisplay string
	Title           string
	Summary         string
	Details         string
	IsActive        bool
	Latitude        float64
	Longitude       float64
	LocationRadius  float64
	Address         string
}

// ConferenceListVMRead
type ConferenceListVMRead struct {
	Base
	ClientID    uuid.UUID
	ClientName  string
	Conferences []models.Conference

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}
