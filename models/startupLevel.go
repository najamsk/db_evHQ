package models

type StartupLevel struct {
	Base
	Name             string
	ClientID         int
	ConferenceID     int
	DisplayLevelSort int
	Startups         []Startup

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}