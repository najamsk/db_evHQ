package models
import(
	"github.com/satori/go.uuid"
)

type Sponsor struct {
	Base
	Name             string
	IsActive         bool
	SponsorLevelID   uuid.UUID
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	SortOrder int
	Description string 
	SocialMedia  SocialMedia `gorm:"embedded"`

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}