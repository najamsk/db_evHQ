package viewmodels

//SubscriptionPayment viewmodel
type SubscriptionPayment struct {
	Base
	IsActive        bool
	Ammount         float64
	AmmountCurrency string
	PaymentLog      string
	PaymentGateway  string
	SubscriptionID  int

	TransactionID string

	// CreatedAt       time.Time
	// UpdatedAt       time.Time
	// DeletedAt       time.Time
}
