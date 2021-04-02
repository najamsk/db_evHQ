package models

import (
	"time"
	uuid "github.com/satori/go.uuid"
	)

type Ticket struct {
	Base
	IsActive        bool
	ValidFrom       time.Time
	ValidTo         time.Time
	ClientID        uuid.UUID
	ConferenceID    uuid.UUID
	SerialNo        string //uuid maybe?
	BookedBy        uuid.UUID    //who purchased this
	IsConsumed      bool
	ConsumedBy      uuid.UUID       //should be same person who booked
	ConsumedAt      time.Time //defualt value should be nil until its consumed
	ConsumedSession int
	SoldBy          uuid.UUID
	TicketTypeID    uuid.UUID

	// CreatedAt       time.Time
	// UpdatedAt       time.Time
	// DeletedAt       time.Time
}