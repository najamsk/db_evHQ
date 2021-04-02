package viewmodels

import (
	"time"

	//"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

// //ConferenceEditVMRead for edit form
type SessionEditVMRead struct {
	Base
	ClientID         uuid.UUID
	ClientName       string
	ConfName         string
	ConfID           uuid.UUID
	StartDate        string
	EndDate          string
	DurationDisplay  string
	Title            string
	Summary          string
	Details          string
	IsActive         bool
	Latitude         float64
	Longitude        float64
	LocationRadius   float64
	Address          string
	Poster           string
	Thumbnail        string
	StartDateDisplay string
	EndDateDisplay   string
	PosterURL        string
	ThumbnailURL	string
	Venue            string
	IsFeatured		bool
	SortOrder		 int
}

// //ConferenceEditVMWrite maybe use this vm from controller to service insteal of too many func params.
// type ConferenceEditVMWrite struct {
// 	ID string
// 	Base
// 	ClientID        uuid.UUID
// 	ClientName      string
// 	StartDate       time.Time
// 	EndDate         time.Time
// 	DurationDisplay string
// 	Title           string
// 	Summary         string
// 	Details         string
// 	IsActive        bool
// 	Latitude        float64
// 	Longitude       float64
// 	LocationRadius  float64
// 	Address         string
// }

// SessionListVMRead should list all sessions
type SessionListRead struct{ 
	Base					// use with count
	StartDate        time.Time
	EndDate          time.Time
	Title			string
	IsActive  		bool
	SpeakerCount    int8
	SortOrder 	int8
}

type SessionListVMRead struct {
	Base
	ClientID   uuid.UUID
	ClientName string
	ConfName   string
	ConfID     uuid.UUID
	Sessions   []SessionListRead
	SpeakerFirstName  string					// used in listing byspeaker
	SpeakerLastName string 	
	SpeakerID  uuid.UUID					// used in listing byspeaker

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}

// SessionListVMRead should list all sessions
type SessionCreateVMWrite struct {
	Base
	ClientID         uuid.UUID
	ClientName       string
	ConfName         string
	ConfID           uuid.UUID
	StartDate        time.Time
	EndDate          time.Time
	DurationDisplay  string
	Title            string
	Summary          string
	Details          string
	IsActive         bool
	Latitude         float64
	Longitude        float64
	LocationRadius   float64
	Address          string
	Poster           string
	Thumbnail        string
	StartDateDisplay string
	EndDateDisplay   string
	Venue            string
	IsFeatured		bool
	SortOrder		 int

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}

// SessionListVMRead should list all sessions
type SessionCreateVMRead struct {
	Base
	ClientID   uuid.UUID
	ClientName string
	ConfName   string
	ConfID     uuid.UUID
	SessionID  uuid.UUID                    // we use it in speaker create

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}
