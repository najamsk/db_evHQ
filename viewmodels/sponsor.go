package viewmodels

import (
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

type SponsorVm struct {
	ID        uuid.UUID
	Name      string
	IsActive  bool
	SortOrder int8
	Type      string
}
type SponsorVmRead struct {
	ClientID     uuid.UUID
	ClientName   string
	ConferenceID uuid.UUID
	SessionID    uuid.UUID
	Sponsors     []SponsorVm
	Sponsorlevel []models.SponsorLevel
}

type SponsorEditVmRead struct {
	Base
	Name              string
	IsActive          bool
	SponsorLevelID    uuid.UUID
	ClientID          uuid.UUID
	ConferenceID      uuid.UUID
	SortOrder         int
	Description       string
	Facebook          string
	Youtube           string
	Twitter           string
	Linkedin          string
	ClientName        string
	Sponsorlevel_ID   uuid.UUID
	SponsorLevel_name string
	Sponsorlevel      []models.SponsorLevel
}

type SponsorEditVmWrite struct {
	Base
	Name         string
	Bio          string
	ClientID     uuid.UUID
	SortOrder    int
	ConferenceID uuid.UUID
	Facebook     string
	Twitter      string
	Linkedin     string
	Youtube      string
	Sponlevel    uuid.UUID
	IsActive bool
}


type SponsorCreateVmWrite struct {
	Name         string
	Bio          string
	ClientID     uuid.UUID
	SortOrder    int
	ConferenceID uuid.UUID
	Facebook     string
	Twitter      string
	Linkedin     string
	Youtube      string
	Sponlevel    uuid.UUID
	IsActive bool
}

