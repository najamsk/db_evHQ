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

// ConferenceService will do dirty work
type ConferenceService struct {
}

// GetConferenceByID should return vm for edit
func (srv *ConferenceService) GetConferenceByID(ConferenceID string) (*viewmodels.ConferenceEditVMRead, error) {
	fmt.Printf("GetConferenceByID func got clientid passed  = %v \n", ConferenceID)

	confID, _ := uuid.FromString(ConferenceID)

	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	imgrepo := repositories.Images{}
	confDB, confError := confrepo.GetByID(confID)

	if confError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", confError)
		return nil, confError

	}

	clientDB, err := clientrepo.GetByID(confDB.ClientID)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}
	// poster, err := imgRepo.GetImage(confDB.ID,"conference", "poster")
	// fmt.Println("poster",poster)
	// if err != nil {
	// 	fmt.Printf("cant find clientdb =  %v \n", err)
	// 	return nil, err

	// }
	hqzone := revel.Config.StringDefault("hq.timezone", "Asia/Karachi")
	loc, _ := time.LoadLocation(hqzone)
	JavascriptISOString := "01/02/2006 15:04:05"

	vm := viewmodels.ConferenceEditVMRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.Title = confDB.Title
	vm.Summary = confDB.Summary
	vm.DurationDisplay = confDB.DurationDisplay
	vm.Details = confDB.Details
	// vm.StartDate = confDB.StartDate
	// vm.EndDate = confDB.EndDate
	vm.IsActive = confDB.IsActive
	vm.ID = confDB.ID

	vm.Address = confDB.Address
	vm.Latitude = confDB.GeoLocation.GeoLocationLat
	vm.Longitude = confDB.GeoLocation.GeoLocationLong
	vm.LocationRadius = confDB.GeoLocation.Radius

	vm.StartDate = confDB.StartDate.In(loc).Format(JavascriptISOString)
	vm.EndDate = confDB.EndDate.In(loc).Format(JavascriptISOString)

	poster, err := imgrepo.GetImage(confDB.ID, "conference", "poster")
	if err == nil {
		vm.PosterURL = poster.BasicURL + poster.ImageURLPrefix + "/" + poster.Name

	}
	thumbnail, err := imgrepo.GetImage(confDB.ID, "conference", "thumbnail")
	if err == nil {
		vm.ThumbnailURL = thumbnail.BasicURL + thumbnail.ImageURLPrefix + "/" + thumbnail.Name

	}
	return &vm, nil
	// return nil, nil
}

// UpdateConference by confid
func (srv *ConferenceService) UpdateConference(ConfData viewmodels.ConferenceEditVMWrite) (uuid.UUID, error) {
	confID, _ := uuid.FromString(ConfData.ID)

	confrepo := repositories.Conferences{}
	confModel, confdberr := confrepo.GetByID(confID)
	if confdberr != nil {
		fmt.Printf("confercnes by id returns error : %v\n", confdberr)
		return uuid.Nil, confdberr

	}
	confModel.Title = ConfData.Title
	confModel.Summary = ConfData.Summary
	confModel.Details = ConfData.Details
	confModel.DurationDisplay = ConfData.DurationDisplay
	confModel.IsActive = ConfData.IsActive
	confModel.StartDate = ConfData.StartDate
	confModel.ClientID = ConfData.ClientID
	confModel.EndDate = ConfData.EndDate
	confModel.Address = ConfData.Address
	confModel.GeoLocation = models.GeoLocation{GeoLocationLat: ConfData.Latitude, GeoLocationLong: ConfData.Longitude, Radius: ConfData.LocationRadius}
	confModel.ID = confID

	confError := confrepo.UpdateConference(&confModel)

	if confError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", confError)
		return uuid.Nil, confError

	}
	return confModel.ID, nil
}

// GetClientNameByClientID by confid
func (srv *ConferenceService) GetClientNameByClientID(ClientID uuid.UUID) string {
	clientrepo := repositories.Clients{}
	clientDB, errDB := clientrepo.GetByID(ClientID)
	if errDB != nil {
		fmt.Printf("erorr while getting client name by id  = %v", errDB)
		return ""
	}
	return clientDB.Name
}

// CreateConference by confid
func (srv *ConferenceService) CreateConference(ConfData viewmodels.ConferenceEditVMWrite) (uuid.UUID, error) {
	// confID, _ := uuid.FromString(ConfData.ID)

	confrepo := repositories.Conferences{}
	// conferencerepo := repositories.Conferences{}
	conference := models.Conference{
		// ID:              confID,
		Title:           ConfData.Title,
		Summary:         ConfData.Summary,
		Details:         ConfData.Details,
		DurationDisplay: ConfData.DurationDisplay,
		IsActive:        ConfData.IsActive,
		StartDate:       ConfData.StartDate,
		ClientID:        ConfData.ClientID,
		EndDate:         ConfData.EndDate,
		Address:         ConfData.Address,
		GeoLocation:     models.GeoLocation{GeoLocationLat: ConfData.Latitude, GeoLocationLong: ConfData.Longitude, Radius: ConfData.LocationRadius},
	}
	// conference.ID = confID

	confError := confrepo.CreateConference(&conference)

	if confError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", confError)
		return uuid.Nil, confError

	}
	return conference.ID, nil
}

//GetConferencesByClientID will return viewmodel
func (srv *ConferenceService) GetConferencesByClientID(ClientID string) (*viewmodels.ConferenceListVMRead, error) {
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
