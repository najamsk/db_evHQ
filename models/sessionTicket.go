package models

type SessionTicket struct {
	Base
	TicketTypeID   int
	TicketIssued   int64 //starts with non zero number
	TicketConsumed int64 //starts with zero number and should reach to ticketissued
	SessionID      int
}