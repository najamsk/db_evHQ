package models

import(
	"github.com/satori/go.uuid")

type Conference_speakers struct {
	ConferenceID	  uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	SortOrder		 int
}