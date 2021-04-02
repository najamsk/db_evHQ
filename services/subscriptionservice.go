package services

import (
	"fmt"
	"time"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// SubscriptionService will do dirty work
type SubscriptionService struct {
}

//GetSubcriptionDetailsByClientID will return viewmodel
func (srv *SubscriptionService) GetSubcriptionDetailsByClientID(ClientID string) (*viewmodels.SubscriptionVMRead, error) {
	fmt.Printf("clientid passed is = %v \n", ClientID)
	var vmRead = viewmodels.SubscriptionVMRead{IsNewSubscription: false}
	clientid, _ := uuid.FromString(ClientID)

	clientrepo := repositories.Clients{}
	subscriptionrepo := repositories.Subscriptions{}

	clientDB, err := clientrepo.GetByID(clientid)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	sub, modelError := subscriptionrepo.GetSubscriptionByClient(clientid, 0)
	if modelError != nil {
		fmt.Printf("repo sub error = %v \n", modelError)
		vmRead.IsNewSubscription = true
		sub = &models.Subscription{Billed: 300.55, BilledCurrency: "usd", StartDate: time.Now().Local(), EndDate: time.Now().Local().Add(time.Hour * 24)}
		// return c.RenderError(modelError)

	} else {
		vmRead.ID = sub.ID
	}

	// fmt.Printf("service sub id is =%v \n", sub.ID)
	vmRead.ClientName = clientDB.Name
	vmRead.ClientID = clientDB.ID
	vmRead.Remarks = sub.Remarks
	vmRead.IsActive = sub.IsActive
	vmRead.Billed = sub.Billed
	vmRead.BilledCurrency = sub.BilledCurrency

	vmRead.StartDate = sub.StartDate
	vmRead.EndDate = sub.EndDate
	vmRead.StartDateDisplay = sub.StartDateDisplay
	vmRead.EndDateDisplay = sub.EndDateDisplay
	vmRead.CreatedAt = sub.CreatedAt

	if len(sub.Payments) > 0 {
		vmRead.PaymentLog = sub.Payments[0].PaymentLog
		fmt.Printf("client found with subscription = %v \n", sub.Payments[0].PaymentLog)
		vmRead.PaymentGateway = sub.Payments[0].PaymentGateway
	}
	hqzone := revel.Config.StringDefault("hq.timezone", "Asia/Karachi")
	loc, _ := time.LoadLocation(hqzone)
	JavascriptISOString := "01/02/2006 15:04:05"

	vmRead.StartTimeISO = sub.StartDate.In(loc).Format(JavascriptISOString)
	vmRead.EndTimeISO = sub.EndDate.In(loc).Format(JavascriptISOString)

	//set timezone,
	now := sub.CreatedAt.In(loc)

	vmRead.CreatedAtISO = now.Format(JavascriptISOString)

	fmt.Printf("clinet name %v \n", clientDB.Name)

	return &vmRead, nil
}
