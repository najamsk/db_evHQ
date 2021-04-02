package viewmodels

import (

	uuid "github.com/satori/go.uuid"
	"time"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
)
type TicketTypeVMWrite struct{
	Base
	//ClientID     uuid.UUID
	ConfID       uuid.UUID
	Title		 string
	Price        float64
	Currency	 string
	IsActive     bool
	Description string

}
type TicketTypelist struct{
	ID uuid.UUID
	Title           string
	IsActive        bool
	ClientID        uuid.UUID
	ConferenceID    uuid.UUID
	Amount          float64
	AmmountCurrency string
	Description string
	ConsumedTicket int 
	TotalTickect int
}
type TicketTypeVMRead struct{
	Base
	ClientID     uuid.UUID
	ConfID       uuid.UUID
	ClientName   string
	TicketType   [] TicketTypelist
}

type TicketTypeVMEdit struct{
	Base
	ID 			 uuid.UUID
	ClientID     uuid.UUID
	ConfID       uuid.UUID
	ClientName   string
	Title		 string
	Price        float64
	Currency	 string
	IsActive     bool
	Description string

}

type TicketTypeVMEditWrite struct{
	Base
	ID 			 uuid.UUID
	Title		 string
	Price        float64
	Currency	 string
	IsActive     bool
	Description string

}
type TicketCreateRead struct{
	ConfId		uuid.UUID
	ClientId    uuid.UUID
	ClientName	string
	TicketsTypes [] models.TicketType
}
type TicketCreateWrite struct{
	ConfId			uuid.UUID
	Title 			string
	StartRange  	int
	EndRange    	int
	TicketTypeId	uuid.UUID
	IsActive		bool
	StartDate       time.Time
	EndDate         time.Time
}

type TicketVMRead struct{
	Base
	ClientID     uuid.UUID
	ConfID       uuid.UUID
	ClientName   string
	Tickets   [] models.Ticket
}
type TicketVMEditRead struct{
	Base
	ID 			 uuid.UUID
	ClientID     uuid.UUID
	ConfID       uuid.UUID
	ClientName   string
	TicketTypeID  uuid.UUID
	StartDate    string
	Title		 string
	EndDate	 	 string
	IsActive     bool
	TicketTypes   []TicketTypelist

}

type TicketVMEditWrite struct{
	Base
	ID 			 uuid.UUID
	ConfID       uuid.UUID
	TicketTypeID  uuid.UUID
	StartDate    time.Time
	EndDate	 	 time.Time
	IsActive     bool
}