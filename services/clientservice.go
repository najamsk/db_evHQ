package services

import (
	"fmt"

	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	uuid "github.com/satori/go.uuid"
)

// ConferenceService will do dirty work
type ClientService struct {
}

func (srv *ClientService) UpdateClientBasicInformation(ID string, Name string, IsActive bool) error {

	clientrepo := repositories.Clients{}

	recID, _ := uuid.FromString(ID)

	dberr := clientrepo.Update(recID, Name, IsActive)
	if dberr != nil {
		fmt.Printf("updating client returns error: %v \n", dberr)
	}

	return dberr

}

// GetConferenceByID should return vm for edit
func (srv *ClientService) GetConferenceByID(RowID string) (*viewmodels.ConferenceEditVMRead, error) {
	fmt.Printf("clientid passed is = %v \n", RowID)

	recID, _ := uuid.FromString(RowID)

	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	// conferencerepo := repositories.Conferences{}

	confDB, confError := confrepo.GetByID(recID)

	if confError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", confError)
		return nil, confError

	}

	clientDB, err := clientrepo.GetByID(recID)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	vm := viewmodels.ConferenceEditVMRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.Title = confDB.Title
	vm.ID = confDB.ID

	return &vm, nil
	// return nil, nil
}

//GetConferencesByClientID will return viewmodel
func (srv *ClientService) GetConferencesByClientID(ClientID string) (*viewmodels.ConferenceListVMRead, error) {
	fmt.Printf("clientid passed is = %v \n", ClientID)

	clientid, _ := uuid.FromString(ClientID)

	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	// conferencerepo := repositories.Conferences{}

	clientDB, err := clientrepo.GetByID(clientid)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	conferences, confError := confrepo.GetByClient(clientid, 0)

	if confError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", confError)
		return nil, confError

	}

	vm := viewmodels.ConferenceListVMRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.Conferences = conferences

	return &vm, nil
}
