package models

type Speaker struct {
	Base
	Name        string
	IsActive    bool
	Poster      string
	Thumbnail   string
	Sessions    []*Session `gorm:"many2many:session_speakers;"`
	ClientID    int
	Conferences []*Conference `gorm:"many2many:speakers_conferences;"`

	DisplayLevelSort int

	// CreatedAt   time.Time
	// UpdatedAt   time.Time
	// DeletedAt   time.Time
}