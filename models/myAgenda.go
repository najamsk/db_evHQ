package models

import "time"

type MyAgenda struct {
	Base
	UserID    int
	Title     string
	StartTime time.Time
	EndTime   time.Time

	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt time.Time
}