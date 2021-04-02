package models
import(
	"github.com/satori/go.uuid"
)

type SponsorLevel struct {
	Base
	Name             string
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	SortOrder		 int
	IsActive		 bool
	//Sponsors         []Sponsor

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}