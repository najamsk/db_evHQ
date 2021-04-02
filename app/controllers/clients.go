package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Clients controller
type Clients struct {
	Admin
}

// Index action: GET
func (c Clients) Index() revel.Result {

	clientrepo := repositories.Clients{}
	// myval := clientrepo.GetDB()
	// fmt.Printf("my baserepo value is = %v \n", myval)
	clients := clientrepo.GetAll(0) //skipping no record
	// for _, client := range clients {

	// 	fmt.Printf("router client loop name = %v\n", client.Name)
	// }
	// mytext := fmt.Sprintf("%v", len(clients))

	// return c.RenderText(mytext)

	return c.Render(clients)
	// return c.RenderTemplate("App/Index.html")
}

// UpdateClient action: UpdateClient
func (c Clients) UpdateClient(id string, name string, isActive bool) revel.Result {

	fmt.Printf("client id = %v \n", id)
	fmt.Printf("client name = %v \n", name)
	fmt.Printf("client isactive = %v \n", isActive)

	var srv = services.ClientService{}

	err := srv.UpdateClientBasicInformation(id, name, isActive)

	fmt.Printf("client service while updating returns error: %v \n", err)
	// /admin/clients/details/ + id
	c.Flash.Success("Client Successfully Updated")
	return c.Redirect("/admin/clients/details/" + id)

}

// CreatePost action: POST
func (c Clients) CreatePost(name string, isActive bool) revel.Result {
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

// List action: GET
func (c Clients) List() revel.Result {

	return c.Render()
}

// Details action: GET
func (c Clients) Details(id string) revel.Result {
	clientrepo := repositories.Clients{}
	clientid, _ := uuid.FromString(id)
	clientDB, err := clientrepo.GetByID(clientid)
	if err != nil {
		log.Fatal("repo didnt locate client by id")
	}
	fmt.Printf("client found with name = %v \n", clientDB.Name)
	return c.Render(clientDB)
}

func (c Clients) Subscription(id string) revel.Result {
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
func (c Clients) SessionsList() revel.Result {
	return c.Render()
}

// TODO: needs to redirect to listing or thank you page. with flash ? also use reverse routes inside html templates
// Create action: Create
func (c Clients) Create(id string) revel.Result {

	fmt.Println("func create called on conferenes controller")
	// clientrepo := repositories.Clients{}
	// clientid, _ := uuid.FromString(id)
	// clientDB, err := clientrepo.GetByID(clientid)
	// if err != nil {
	// 	log.Fatal("repo didnt locate client by id")
	// }
	// fmt.Printf("client found with name = %v \n", clientDB.Name)

	return c.Render()
}

// func (c Clients) Create(name string, company string) revel.Result {
// 	//its working using c.Params.Get for post values as well
// 	// name := c.Params.Get("name")
// 	fmt.Printf("name is : %v", company)

// 	greet := company + " thank you."
// 	return c.RenderText(greet)
// }

// Payments action: Payments
func (c Clients) Payments() revel.Result {
	return c.Render()
}

func (c Clients) Search(name string) revel.Result {
	clientrepo := repositories.Clients{}
	clientName := strings.ToLower(name)
	clients, Dberr := clientrepo.SearchByName(clientName, 0) //skipping no record
	if Dberr != nil {
		fmt.Println("sreach client shows error", Dberr)
		return c.Redirect("/admin/clients")
	}
	return c.Render(clients)
}
