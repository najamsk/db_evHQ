package services

import (
	"fmt"
	"time"
	"strconv"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"

	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// ConferenceService will do dirty work
type TicketService struct {
}


func (srv *TicketService) GetByconfID(confid uuid.UUID) (*viewmodels.TicketVMRead, error) {
	fmt.Printf("confid passed in ticket service at getbyconfid", confid)
	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	ticktrepo := repositories.Tickets{}
	confDb, err := confrepo.GetByID(confid)
	if err != nil {
		fmt.Println("get conference by id err at ticket service", err)
		return nil, err

	}
	clientDB, err := clientrepo.GetByID(confDb.ClientID)
	if err != nil {
		fmt.Println("get client by id err at ticket service", err)
		return nil, err

	}
	ticketDB, err := ticktrepo.GetByConference(confid)
	if err != nil {
		fmt.Println("get ticket by confid err at ticket service", err)
		return nil, err

	}
	vm := viewmodels.TicketVMRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.ConfID = confDb.ID
	vm.Tickets = ticketDB
	return &vm, nil
}

func (srv *TicketService) CreateTiket(obj viewmodels.TicketCreateWrite)error {
	confrepo := repositories.Conferences{}
	ticktrepo := repositories.Tickets{}
	//title:=obj.Title
	fmt.Println("confiddd passed in ticket service", obj.ConfId)
	confDB, repoErr := confrepo.GetByID(obj.ConfId)
	if repoErr != nil {
		fmt.Println("cant find confdb by id in ticketservice")
		return repoErr
	}

	 for i := obj.StartRange;  i<=obj.EndRange; i++  {
		TicketData := models.Ticket{
			IsActive:   obj.IsActive,
			ConferenceID:   obj.ConfId,
			ClientID:	confDB.ClientID,
			SerialNo    :obj.Title+strconv.Itoa(i),
			TicketTypeID:   obj.TicketTypeId,
			ValidFrom    :obj.StartDate,
			ValidTo: obj.EndDate,
		}
		TicketErr := ticktrepo.CreateTicket(&TicketData)
		if TicketErr != nil {
			fmt.Println( TicketData.SerialNo+"cant create ticket type ticketservice")
			return TicketErr
		}

     }


	return nil
}

func (srv *TicketService) GetTicketById(id string, confid string) (*viewmodels.TicketVMEditRead, error) {
	hqzone := revel.Config.StringDefault("hq.timezone", "Asia/Karachi")
	loc, _ := time.LoadLocation(hqzone)
	JavascriptISOString := "01/02/2006 15:04"
	
	TicketID, err := uuid.FromString(id)
	if err != nil {
		fmt.Println("ticket type id show error in conversion at GetTicketTypeById", err)
		return nil, err
	}
	confidID, conferr := uuid.FromString(confid)
	if conferr != nil {
		fmt.Println("confid show error in conversion at GetTicketTypeById", conferr)
		return nil, conferr
	}
	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	tiktTyperepo := repositories.TicketsType{}
	ticketrepo:= repositories.Tickets{}
	// conferencerepo := repositories.Conferences{}

	confDB, repoError := confrepo.GetByID(confidID)
	if repoError != nil {
		fmt.Printf("conference by id not found with error in ticket service : %v\n", repoError)
		return nil, repoError

	}
	clientDB, clientError := clientrepo.GetByID(confDB.ClientID)
	if repoError != nil {
		fmt.Printf("client by id not found with error in ticket service  : %v\n", clientError)
		return nil, clientError

	}
	ticketdb, resultError := ticketrepo.GetById(TicketID)
	if resultError != nil {
		fmt.Println("GetTypeById show error in ticket service", resultError)
		return nil, resultError
	}
	TickettypeDB, clientError := tiktTyperepo.GetByConference(confidID)
	if repoError != nil {
		fmt.Printf("client by id not found with error in ticket service  : %v\n", clientError)
		return nil, clientError

	}

	vm := viewmodels.TicketVMEditRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.ConfID = confDB.ID
	vm.ID = TicketID
	vm.TicketTypeID = ticketdb.TicketTypeID
	vm.StartDate = ticketdb.ValidFrom.In(loc).Format(JavascriptISOString)
	vm.Title=ticketdb.SerialNo
	vm.EndDate = ticketdb.ValidTo.In(loc).Format(JavascriptISOString)
	vm.IsActive = ticketdb.IsActive
	vm.TicketTypes=TickettypeDB
	return &vm, nil
}
func(srv *TicketService) UpdateTicket(obj viewmodels.TicketVMEditWrite) error{
	confrepo := repositories.Conferences{}
	ticktrepo := repositories.Tickets{}
	
	confDB, repoErr := confrepo.GetByID(obj.ConfID)
	if repoErr != nil {
		fmt.Println("cant find confdb by id in ticketservice")
		return repoErr
	}
	
	TicketModel,ERR:=ticktrepo.GetById(obj.ID)
	if ERR != nil {
		fmt.Println( "cant get ticketbyid ticketservice")
		return ERR
	}

	TicketModel.IsActive=   obj.IsActive
	TicketModel.ConferenceID=   obj.ConfID
	TicketModel.ClientID=	confDB.ClientID
	TicketModel.TicketTypeID=   obj.TicketTypeID
	TicketModel.ValidFrom=    obj.StartDate
	TicketModel.ValidTo= obj.EndDate
	TicketModel.ID=obj.ID

	
	TicketErr := ticktrepo.UpdateTicket(TicketModel)
	if TicketErr != nil {
		fmt.Println( "cant update tickets ticketservice")
		return TicketErr
	}


	return nil
}
