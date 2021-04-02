package controllers

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Subscriptions controller
type Subscriptions struct {
	//*revel.Controller
	Admin
}

// Index action: GET
func (c Subscriptions) Index(id string) revel.Result {

	clientrepo := repositories.Clients{}
	// myval := clientrepo.GetDB()
	// fmt.Printf("my baserepo value is = %v \n", myval)
	clients := clientrepo.GetAll(0) //skipping no record
	for _, client := range clients {

		fmt.Printf("router client loop name = %v\n", reflect.TypeOf(client.ID))
	}
	
	return c.RenderHTML("Subscriptions/SubscriptionDetail.html")
	// return c.RenderText(id)
}

//ByClientID action
func (c Subscriptions) ByClientID(id string) revel.Result {

	var srv = services.SubscriptionService{}
	subscriptionvm, vmerror := srv.GetSubcriptionDetailsByClientID(id)
	if vmerror != nil {
		fmt.Println("subscription service returns errors")

	}
	fmt.Printf("subvm.clientid  = %v\n", subscriptionvm.ClientID)
	fmt.Printf("subvm  = %v\n", subscriptionvm)

	c.ViewArgs["subscription"] = subscriptionvm

	return c.RenderTemplate("Subscriptions/SubscriptionDetail.html")

}

//BySubscriptionID action
func (c Subscriptions) BySubscriptionID(id string) revel.Result {

	// clientrepo := repositories.Clients{}
	// // myval := clientrepo.GetDB()
	// // fmt.Printf("my baserepo value is = %v \n", myval)
	// clients := clientrepo.GetAll(0) //skipping no record
	// for _, client := range clients {

	// 	fmt.Printf("router client loop name = %v\n", reflect.TypeOf(client.ID))
	// }
	return c.RenderText("by subsctiption: " + id)
}

// CreatePost action: POST
func (c Subscriptions) CreatePost(name string, isActive bool) revel.Result {
	// revel.Config.StringDefault("hq.developer.name", "nsa")
	client := models.Client{Name: name, IsActive: isActive}

	clientrepo := repositories.Clients{}

	clientid, err := clientrepo.Create(client)

	if err != nil {
		log.Fatal("repo throws error creating new client")
		return c.RenderText("sorry can't create this client now.")
	}

	fmt.Printf("name: %v, isActive:%v, Id= %v \n", client.Name, client.IsActive, clientid)
	resultMsg := client.Name + " created"
	c.ViewArgs["thanks"] = resultMsg
	c.Flash.Success(resultMsg)

	// return c.RenderTemplate("Clients/Index.html")
	return c.Redirect(Clients.Index)
}

// CreatePost action: POST
func (c Subscriptions) UpdateSubscription(
	clientID string,
	subscriptionid string,
	startDate string,
	startDateDisplay string,
	endDate string,
	endDateDisplay string,
	billed float64,
	billedCurrency string,
	paymentGateway string,
	remarks string,
	paymentLog string,
	isActive bool) revel.Result {

	// just printing all the values from Form to console
	for k, v := range c.Params.Values {
		fmt.Printf("param %s is = %s \n", k, v)
	}
	//working on update subscription func.
	//it should first get active sub based on client id.
	fmt.Println("updatesubscription func activated. should run vlaidations.")
	fmt.Printf("subscriptionid from hidden field is = %v", subscriptionid)

	if len(subscriptionid) != 0 {
		fmt.Println("we should create new subscription from client")
	}

	clientid, _ := uuid.FromString(clientID)

	i, err := strconv.ParseInt(startDate, 10, 64)
	if err != nil {
		panic(err) //maybe shourld return flash error to form page
	}

	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Printf("start time from timestamp = %v \n", time.Unix(i, 0))
	fmt.Printf("start time from timestamp/1000 = %v \n", time.Unix(i/1000, 0))
	fmt.Printf("start time from timestamp/1000.utc = %v \n", time.Unix(i/1000, 0).UTC())
	startDateStamp := time.Unix(i/1000, 0).UTC()

	i, err = strconv.ParseInt(endDate, 10, 64)
	if err != nil {
		panic(err) //maybe shourld return flash error to form page
	}
	endDateStamp := time.Unix(i/1000, 0).UTC()

	fmt.Println("utc startdate is ...")
	fmt.Println(startDateStamp) //should return utc datetime
	fmt.Println(endDateStamp)   //should return utc datetime
	fmt.Printf("billed = %v", billed)

	sub := models.Subscription{IsActive: isActive,
		StartDate:        startDateStamp,
		EndDate:          endDateStamp,
		StartDateDisplay: startDateDisplay,
		EndDateDisplay:   endDateDisplay,
		DurationDisplay:  "", //should be send on differance between start and end
		ClientID:         clientid,
		Billed:           billed,
		BilledCurrency:   billedCurrency,
		Remarks:          remarks,
	}
	subscriptionRepo := repositories.Subscriptions{}
	if len(subscriptionid) == 0 {
		fmt.Println("we should create new subscription from client")

		sub.Payments = []*models.SubscriptionPayment{&models.SubscriptionPayment{
			IsActive:        isActive,
			Ammount:         billed,
			AmmountCurrency: billedCurrency,
			PaymentLog:      paymentLog,
			PaymentGateway:  paymentGateway,
		},
		}

		subid, errdb := subscriptionRepo.Create(&sub)
		if errdb != nil {
			log.Fatal(errdb)
		}
		fmt.Printf("sub is = %v", sub)

		fmt.Printf("sub created with id = %v", subid)

	} else {
		frmSubID, errSubID := uuid.FromString(subscriptionid)
		if errSubID != nil {
			fmt.Printf("cant parse subscid %v", errSubID)
		}
		fmt.Printf("subscription id from Form post is =  %v", frmSubID)
		sub.ID = frmSubID
		updated := subscriptionRepo.UpdateSubscription(&sub)

		if updated == false {
			fmt.Println("cant update subscription")
		}

	}
	redirecturl := "/admin/subscriptions/client/" + clientID
	fmt.Printf("about to redirect to : %v \n", redirecturl)
	// revel.ReverseURL("Subscriptions.ByClientID", clientID)
	return c.Redirect(redirecturl)
	// return c.RenderTemplate("Clients/Index.html")
	// return c.RenderTemplate("Subscriptions/SubscriptionDetail.html")
	// return c.Redirect(Clients.Index)
}

// List action: GET
func (c Subscriptions) List() revel.Result {
	return c.Render()
}

// Details action: GET
func (c Subscriptions) Details(id string) revel.Result {
	clientrepo := repositories.Clients{}
	clientid, _ := uuid.FromString(id)
	clientDB, err := clientrepo.GetByID(clientid)
	if err != nil {
		log.Fatal("repo didnt locate client by id")
	}
	fmt.Printf("client found with name = %v \n", clientDB.Name)

	return c.Render(clientDB)
}

func (c Subscriptions) Subscription(id string) revel.Result {
	clientrepo := repositories.Clients{}
	clientid, _ := uuid.FromString(id)
	clientDB, err := clientrepo.GetByID(clientid)
	if err != nil {
		log.Fatal("repo didnt locate client by id")
	}
	fmt.Printf("client found with name = %v \n", clientDB.Name)
	return c.Render(clientDB)
}

// SessionsList action: GET
func (c Subscriptions) SessionsList() revel.Result {
	return c.Render()
}

// TODO: needs to redirect to listing or thank you page. with flash ? also use reverse routes inside html templates
// Create action: Create
func (c Subscriptions) Create() revel.Result {

	return c.Render()
}

// Payments action: Payments
func (c Subscriptions) Payments() revel.Result {
	return c.Render()
}
