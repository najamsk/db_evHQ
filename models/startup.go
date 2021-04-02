package models

type Startup struct {
	Base
	Name             string
	IsActive         bool
	Poster           string
	Thumbnail        string
	StartupLevelID   int
	ClientID         int
	ConferenceID     int
	DisplayLevelSort int

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}