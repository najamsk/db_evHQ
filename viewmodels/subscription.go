package viewmodels

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// // Subscription viewmodel to work with ui
// type Subscription struct {
// 	Base
// 	IsActive         bool
// 	StartDate        time.Time
// 	StartDateDisplay string
// 	EndDate          time.Time
// 	EndDateDisplay   string
// 	DurationDisplay  string
// 	ClientID         int
// 	ClientName       string
// 	Billed           float64
// 	BilledCurrency   string
// 	Payments         []SubscriptionPayment
// 	Remarks          string
// 	// CreatedAt        time.Time
// 	// UpdatedAt        time.Time
// 	// DeletedAt        time.Time
// }

type SubscriptionVMRead struct {
	Base
	IsActive          bool
	StartDate         time.Time
	StartDateDisplay  string
	EndDate           time.Time
	EndDateDisplay    string
	DurationDisplay   string
	Billed            float64
	BilledCurrency    string
	Payments          []SubscriptionPayment
	Remarks           string
	PaymentLog        string
	IsNewSubscription bool
	ClientID          uuid.UUID
	ClientName        string
	StartTimeISO      string
	EndTimeISO        string
	PaymentGateway    string
	CreatedAtISO      string
	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}
