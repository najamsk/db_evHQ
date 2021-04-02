package viewmodels

import (
	uuid "github.com/satori/go.uuid"
)

type SpeakerVM struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Organization string
	Designation string
	SortOrder	string
	IsActive   bool
}

type SpeakerListVMRead struct {
	Base
	ClientID     uuid.UUID
	ClientName   string
	ConferenceID uuid.UUID
	Speakers     []SpeakerVM
	SessionID    uuid.UUID
	SessionTitle string

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}

type SpeakerEditVMRead struct {
	Base
	ID               uuid.UUID
	ClientID         uuid.UUID
	ConfID           uuid.UUID
	ClientName       string
	FirstName        string
	LastName         string
	PhoneNumber      string
	Bio              string
	Organization     string
	Email            string
	Designation      string
	SessionID        uuid.UUID
	SessionWeight    int
	ConferenceWeight int
	ProfileURL       string
	PosterURL        string
	Facebook         string
	Youtube          string
	Linkedin         string
	Twitter          string
}
type SpeakerEditVMWrite struct {
	//ID string
	Base
	ClientID         uuid.UUID
	FirstName        string
	LastName         string
	Email            string
	Organization     string
	Designation      string
	PhoneNumber      string
	Bio              string
	SessionWeight    int
	ConferenceWeight int
	SessionID        uuid.UUID
	ConferenceID     uuid.UUID
	Facebook         string
	Youtube          string
	Linkedin         string
	Twitter          string
}
type AddSessionSpeakerVMWrite struct {
	//ID string
	Base
	FirstName        string
	LastName         string
	Email            string
	Organization     string
	Designation      string
	PhoneNumber      string
	Bio              string
	ConfID           uuid.UUID
	SessionID        uuid.UUID
	UserId           string
	SessionWeight    int
	ConferenceWeight int
	Facebook         string
	Youtube          string
	Linkedin         string
	Twitter          string
}
