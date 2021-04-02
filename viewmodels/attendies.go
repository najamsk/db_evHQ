package viewmodels

import (
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

type AttendiesListVMRead struct {
	Base
	ClientID     uuid.UUID
	ClientName   string
	ConferenceID uuid.UUID
	Attendies    []models.User
}
type AttendiesEditVMRead struct {
	Base
	ID		   uuid.UUID
	ClientID   uuid.UUID
	ConfID     uuid.UUID
	ClientName string
	// Attendies  *models.User
	FirstName    string
	LastName     string
	PhoneNumber  string
	Bio          string
	Organization string
	Email        string
	Designation  string
}

//AttendiesEditVMWrite
type AttendiesEditVMWrite struct {
	Base
	ClientID     uuid.UUID
	ConfID       uuid.UUID
	ClientName   string
	FirstName    string
	LastName     string
	PhoneNumber  string
	Bio          string
	Organization string
	Email        string
	Designation  string
}
