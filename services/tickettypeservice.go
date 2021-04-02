package services

import (
	"fmt"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"

	//"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// ConferenceService will do dirty work
type TicketTypeService struct {
}

func (srv *TicketTypeService) GetTicketTypeById(id string, confid string) (*viewmodels.TicketTypeVMEdit, error) {
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
	ticketrepo := repositories.TicketsType{}
	//imgrepo:= repositories.Images{}
	// conferencerepo := repositories.Conferences{}

	result, resultError := ticketrepo.GetTypeById(TicketID)
	if resultError != nil {
		fmt.Println("GetTypeById show error in ticket service", resultError)
		return nil, resultError
	}
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
	vm := viewmodels.TicketTypeVMEdit{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.ConfID = confDB.ID
	vm.ID = TicketID
	vm.Title = result.Title
	vm.Price = result.Amount
	vm.Currency = result.AmmountCurrency
	vm.IsActive = result.IsActive
	vm.Description=result.Description
	return &vm, nil
}
func (srv *TicketTypeService) GetByconfID(confid uuid.UUID) (*viewmodels.TicketTypeVMRead, error) {
	fmt.Printf("confid passed in ticket service at getbyconfid", confid)
	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	ticktrepo := repositories.TicketsType{}
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
	vm := viewmodels.TicketTypeVMRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.ConfID = confDb.ID
	vm.TicketType = ticketDB
	return &vm, nil
}
func (srv *TicketTypeService) CreateTiketType(obj viewmodels.TicketTypeVMWrite) (uuid.UUID,error) {
	confrepo := repositories.Conferences{}
	ticktrepo := repositories.TicketsType{}
	fmt.Println("confid passed in ticket service", obj.ConfID)
	confDB, repoErr := confrepo.GetByID(obj.ConfID)
	if repoErr != nil {
		fmt.Println("cant find confdb by id in ticketservice")
		return uuid.Nil,repoErr
	}

	TicketType := models.TicketType{
		Title:           obj.Title,
		IsActive:        obj.IsActive,
		ClientID:        confDB.ClientID,
		ConferenceID:    obj.ConfID,
		Amount:          obj.Price,
		AmmountCurrency: obj.Currency,
		Description:obj.Description,
	}
	TicketErr := ticktrepo.CreateTicketType(&TicketType)
	if TicketErr != nil {
		fmt.Println("cant create ticket type ticketservice")
		return uuid.Nil,TicketErr
	}
	return TicketType.ID,nil
}
func (srv *TicketTypeService) UpdateTiketType(obj viewmodels.TicketTypeVMEditWrite) error {
	fmt.Printf("about to update session in side service. \n")
	fmt.Printf("tickettype passed1 = %v \n", obj.ID)

	Ticketrepo := repositories.TicketsType{}
	TicketTypeModel,ERR:=Ticketrepo.GetTypeById(obj.ID)
	if ERR != nil {
		fmt.Println("UpdateTiketType get by id show error at ticket service", ERR)
		return ERR
	}
	TicketTypeModel.Title=           obj.Title
	TicketTypeModel.IsActive=        obj.IsActive
	TicketTypeModel.Amount=          obj.Price
	TicketTypeModel.AmmountCurrency= obj.Currency
	TicketTypeModel.ID=obj.ID
	TicketTypeModel.Description=obj.Description
	repoErr := Ticketrepo.UpdateTiketType(TicketTypeModel)
	if repoErr != nil {
		fmt.Println("UpdateTiketType show error at ticket service", repoErr)
		return repoErr
	}
	return nil
}