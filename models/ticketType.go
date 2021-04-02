package models
import(
	"github.com/satori/go.uuid"
)

type TicketType struct {
	Base
	Title           string
	IsActive        bool
	ClientID        uuid.UUID
	ConferenceID    uuid.UUID
	Amount          float64
	AmmountCurrency string
	Description string
	//Tickets         []Ticket
	
	// CreatedAt       time.Time
	// UpdatedAt       time.Time
	// DeletedAt       time.Time
}