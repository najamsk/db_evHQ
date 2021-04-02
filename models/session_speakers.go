package models
import(
	"github.com/satori/go.uuid")

type Session_speakers struct {
	SessionID	  	  uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	SortOrder		 int
}