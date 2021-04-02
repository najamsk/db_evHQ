package controllers

import (
	"fmt"
	"strconv"
	"time"
	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Conference controller
type Tickets struct {
	Admin
}

// Index action: GET
func (c Tickets) Index() revel.Result {
	return c.Render()
}

func (c Tickets) Create(confid string) revel.Result {
	confId, _ := uuid.FromString(confid)
	tktSrv := services.TicketTypeService{}
	vm, srvError := tktSrv.GetByconfID(confId)
	if srvError != nil {
		fmt.Printf("tktSrv returns error %v \n", srvError)
		return c.RenderError(srvError)
	}
	
	return c.Render(vm)
}
func (c Tickets) CreatePost(confid string,
	isrange bool,
	end string,
	start string,
	name string,
	tickettype string,
	startDate string,
	endDate string,
	IsActive bool) revel.Result {
	fmt.Println("confid at Createpost", confid)
	fmt.Println("start range at Createpost", start)
	fmt.Println("end range at createpost", end)
	fmt.Println("confid at tickettypeid", tickettype)
	confId, _ := uuid.FromString(confid)
	tikttypeid, _ := uuid.FromString(tickettype)
	var rangeStart int
	var rangeEnd int
	var err error
	if isrange == true {
		rangeStart, err = strconv.Atoi(start)
		if err != nil {
			fmt.Println(" rangeStart is not an integer.")
			c.Flash.Error("Enter valid range")
			return c.Redirect("/admin/conference/tickets/create/" + confid)
		}

		rangeEnd, err = strconv.Atoi(end)
		if err != nil {
			fmt.Println("rangeEnd is not an integer.")
			c.Flash.Error("Enter valid range")
			return c.Redirect("/admin/conference/tickets/create/" + confid)
		}

		TiktRange := rangeEnd - rangeStart
		if TiktRange > 500 {
			fmt.Println("Limit exceed tiket range should be less than 500 ")
			c.Flash.Error("Tiket range should be less than 500")
			return c.Redirect("/admin/conference/tickets/create/" + confid)
		}
	}
	if isrange == false {
		rangeStart=1
		rangeEnd=1

	}
	i, err := strconv.ParseInt(startDate, 10, 64)
	if err != nil {
		panic(err) //maybe shourld return flash error to form page
	}

	startDateStamp := time.Unix(i/1000, 0).UTC()

	i, err = strconv.ParseInt(endDate, 10, 64)
	if err != nil {
		panic(err) //maybe shourld return flash error to form page
	}
	endDateStamp := time.Unix(i/1000, 0).UTC()
	tktSrv := services.TicketService{}
	ticktData := viewmodels.TicketCreateWrite{
		ConfId:       confId,
		Title:        name,
		StartRange:   rangeStart,
		EndRange:     rangeEnd,
		TicketTypeId: tikttypeid,
		IsActive:     IsActive,
		StartDate:    startDateStamp,
		EndDate:      endDateStamp,
	}
	srvError := tktSrv.CreateTiket(ticktData)
	if srvError != nil {
		fmt.Printf("tktSrv returns error %v \n", srvError)
		c.Flash.Error("Serial number Already Exist")
		return c.Redirect("/admin/conference/tickets/create/" + confid)
	}
	c.Flash.Success("Tickect successfuly created")
	return c.Redirect("/admin/conference/tickets/list/" + confid)
}
func (c Tickets) List(Id string) revel.Result {
	confId, _ := uuid.FromString(Id)
	tktSrv := services.TicketService{}
	vm, srvError := tktSrv.GetByconfID(confId)
	if srvError != nil {
		fmt.Printf("tktSrv returns error %v \n", srvError)
		return c.RenderError(srvError)
	}
	return c.Render(vm)
}
func (c Tickets) Edit(confid string, id string) revel.Result {
	var srv = services.TicketService{}

	vm, srvError := srv.GetTicketById(id, confid)
	if srvError != nil {
		fmt.Println("GetTicketTypeById error in ticket controller")
		return c.RenderText(srvError.Error())
	}

	return c.Render(vm)
}
func (c Tickets) EditPost(
	confid string,
	id string,
	end string,
	start string,
	tickettype string,
	startDate string,
	endDate string,
	IsActive bool) revel.Result {

	fmt.Println("edit post func activated")
	fmt.Printf("form confid = %v \n", id)

	var srv = services.TicketService{}
	TicketID, _ := uuid.FromString(id)
	confId, _ := uuid.FromString(confid)
	tikttypeid, _ := uuid.FromString(tickettype)
	i, err := strconv.ParseInt(startDate, 10, 64)
	if err != nil {
		panic(err) //maybe shourld return flash error to form page
	}

	startDateStamp := time.Unix(i/1000, 0).UTC()

	i, err = strconv.ParseInt(endDate, 10, 64)
	if err != nil {
		panic(err) //maybe shourld return flash error to form page
	}
	endDateStamp := time.Unix(i/1000, 0).UTC()

	ticktTypeData := viewmodels.TicketVMEditWrite{
		ID:           TicketID,
		ConfID:       confId,
		TicketTypeID: tikttypeid,
		StartDate:    startDateStamp,
		EndDate:      endDateStamp,
		IsActive:     IsActive,
	}

	srvError := srv.UpdateTicket(ticktTypeData)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderText(srvError.Error())
	}

	// // return c.RenderText("will post data")
	c.Flash.Success("Ticket Successfully Updated")

	return c.Redirect("/admin/conference/" + confid + "/tickets/details/" + id)

}
