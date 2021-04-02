package controllers

import (
	"fmt"
	"strconv"
	//"time"

	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Conference controller
type TicketsType struct {
	Admin
}

// Index action: GET
func (c TicketsType) Index() revel.Result {

	return c.Render()
}

func (c TicketsType) CreateTiketTypePost(
	confid string,
	price string,
	currency string,
	isactive bool,
	desc string,
	name string) revel.Result {

	fmt.Println("edit post func activated")

	fmt.Printf("form title = %v \n", name)
	fmt.Printf("form confid = %v \n", confid)
	fmt.Printf("form price = %v \n", price)
	fmt.Printf("form currency = %v \n", currency)
	fmt.Printf("form isactive = %v \n", isactive)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	var srv = services.TicketTypeService{}

	ConfID, _ := uuid.FromString(confid)
	amount, _ := strconv.ParseFloat(price, 64)
	ticktTypeData := viewmodels.TicketTypeVMWrite{
		ConfID:   ConfID,
		Title:    name,
		Price:    amount,
		Currency: currency,
		IsActive: isactive,
		Description:desc,
	}

	ticktTypeID,srvError := srv.CreateTiketType(ticktTypeData)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderText(srvError.Error())
	}
	fmt.Println(ticktTypeID)
	// // return c.RenderText("will post data")
	// c.Flash.Success("Conference Successfully Updated")

	return c.Redirect("/admin/conference/tickets/create/"+confid)

}

func (c TicketsType) CreateTiketType(confid string) revel.Result {
	var srv = services.ConferenceService{}

	vm, srvError := srv.GetConferenceByID(confid)
	if srvError != nil {
		return c.RenderText(srvError.Error())
	}

	return c.Render(vm)
}
func (c TicketsType) List(Id string) revel.Result {
	confId, _ := uuid.FromString(Id)
	tktSrv := services.TicketTypeService{}
	vm, srvError := tktSrv.GetByconfID(confId)
	if srvError != nil {
		fmt.Printf("tktSrv returns error %v \n", srvError)
		return c.RenderError(srvError)
	}

	return c.Render(vm)
}

func (c TicketsType) EditTiketType(confid string, id string) revel.Result {
	var srv = services.TicketTypeService{}

	vm, srvError := srv.GetTicketTypeById(id, confid)
	if srvError != nil {
		fmt.Println("GetTicketTypeById error in ticket controller")
		return c.RenderText(srvError.Error())
	}
	return c.Render(vm)
}

func (c TicketsType) EditTiketTypePost(
	confid string,
	id string,
	price string,
	currency string,
	isActive bool,
	name string,
	desc string) revel.Result {

	fmt.Println("edit post func activated")

	fmt.Printf("form title = %v \n", name)
	fmt.Printf("form confid = %v \n", confid)
	fmt.Printf("form price = %v \n", price)
	fmt.Printf("form currency = %v \n", currency)
	fmt.Printf("form isactive1 = %v \n", isActive)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	var srv = services.TicketTypeService{}
	TicketID, _ := uuid.FromString(id)
	amount, _ := strconv.ParseFloat(price, 64)
	ticktTypeData := viewmodels.TicketTypeVMEditWrite{
		Title:    name,
		Price:    amount,
		Currency: currency,
		IsActive: isActive,
		ID: TicketID,
		Description:desc,
	}

	srvError := srv.UpdateTiketType(ticktTypeData)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderText(srvError.Error())
	}

	// // return c.RenderText("will post data")
	// c.Flash.Success("Conference Successfully Updated")

	return c.Redirect("/admin/conference/"+confid+"/tickets/type/details/"+id)

}
