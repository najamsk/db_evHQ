package models

type Ambassador struct {
	Base
	Title             string
	ManagerID         int //manager id from users table
	AmbassadorGroupID int
	UserID            int
	IsActive          bool
	Conferences      []*Conference `gorm:"many2many:organizer_conferences;"`
	ClientID         int
	Sessions         []*Session `gorm:"many2many:session_ambassador;"`
	DisplayLevelSort int
	DisplayInList    bool
	DisplayTitle     string

	// CreatedAt         time.Time
	// UpdatedAt         time.Time
	// DeletedAt         time.Time
}