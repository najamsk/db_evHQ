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
type SessionService struct {
}

// GetConferenceByID should return vm for edit
func (srv *SessionService) GetCByID(ID string) (*viewmodels.SessionEditVMRead, error) {
	fmt.Printf("sesion id passed  = %v \n", ID)

	id, _ := uuid.FromString(ID)

	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	sessionrepo := repositories.Sessions{}
	imgrepo := repositories.Images{}
	// conferencerepo := repositories.Conferences{}

	result, resultError := sessionrepo.GetByID(id)

	if resultError != nil {
		fmt.Printf("session by id not found with error  : %v\n", resultError)
		return nil, resultError

	}

	confDB, confError := confrepo.GetByID(result.ConferenceID)
	if confError != nil {
		fmt.Printf("conference by id not found with error  : %v\n", confError)
		return nil, resultError
	}

	fmt.Printf("clien it is =  %v \n", confDB.ClientID)

	clientDB, clientError := clientrepo.GetByID(confDB.ClientID)

	if clientError != nil {
		fmt.Printf("client by id not found with error  =  %v \n", clientError)
		return nil, clientError

	}

	hqzone := revel.Config.StringDefault("hq.timezone", "Asia/Karachi")
	loc, _ := time.LoadLocation(hqzone)
	JavascriptISOString := "01/02/2006 15:04:05"

	vm := viewmodels.SessionEditVMRead{}
	vm.ClientID = clientDB.ID
	vm.ConfID = confDB.ID
	vm.ClientName = clientDB.Name
	vm.ConfName = confDB.Title
	vm.Title = result.Title
	vm.Summary = result.Summary
	vm.DurationDisplay = confDB.DurationDisplay
	vm.Details = result.Details
	// vm.StartDate = confDB.StartDate
	// vm.EndDate = confDB.EndDate
	vm.IsActive = result.IsActive
	vm.ID = result.ID
	vm.Venue = result.Venue
	vm.Address = result.Address
	vm.Latitude = result.GeoLocation.GeoLocationLat
	vm.Longitude = result.GeoLocation.GeoLocationLong
	vm.LocationRadius = result.GeoLocation.Radius
	vm.IsFeatured = result.IsFeatured
	vm.SortOrder = result.SortOrder
	vm.StartDate = result.StartDate.In(loc).Format(JavascriptISOString)
	vm.EndDate = result.EndDate.In(loc).Format(JavascriptISOString)

	poster, err := imgrepo.GetImage(result.ID, "session", "poster")
	if err == nil {
		vm.PosterURL = poster.BasicURL + poster.ImageURLPrefix + "/" + poster.Name

	}
	thumbnail, err := imgrepo.GetImage(result.ID, "session", "thumbnail")
	if err == nil {
		vm.ThumbnailURL = thumbnail.BasicURL + thumbnail.ImageURLPrefix + "/" + thumbnail.Name

	}

	return &vm, nil
	// return nil, nil
}

// GetClientNameByClientID by confid
func (srv *SessionService) GetClientNameByClientID(ClientID uuid.UUID) string {
	clientrepo := repositories.Clients{}
	clientDB, errDB := clientrepo.GetByID(ClientID)
	if errDB != nil {
		fmt.Printf("erorr while getting client name by id  = %v", errDB)
		return ""
	}
	return clientDB.Name
}

// CreateConference by confid
func (srv *SessionService) CreateConference(ConfData viewmodels.ConferenceEditVMWrite) error {
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
		return confError

	}
	return nil
}

//GetConferencesByClientID will return viewmodel
func (srv *SessionService) GetConferencesByClientID(ClientID string) (*viewmodels.ConferenceListVMRead, error) {
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

//GetConferencesByClientID will return viewmodel
func (srv *SessionService) GetSessionsByConferenceID(ConferenceID string) (*viewmodels.SessionListVMRead, error) {
	fmt.Printf("confid passed  in session GetSessionsByConferenceID is = %v \n", ConferenceID)

	confid, _ := uuid.FromString(ConferenceID)

	// clientrepo := repositories.Clients{}
	// confrepo := repositories.Conferences{}
	sessionRepo := repositories.Sessions{}
	confService := ConferenceService{}
	// conferencerepo := repositories.Conferences{}

	// clientDB, err := clientrepo.GetByID(confid)

	// if err != nil {
	// 	fmt.Printf("cant find clientdb =  %v \n", err)
	// 	return nil, err

	// }

	confDB, confError := confService.GetConferenceByID(ConferenceID)
	if confError != nil {
		fmt.Printf("conferece cant be found by id =  %v \n", confid)
	}

	sessions, sessionErr := sessionRepo.GetByConference(confid, 0)

	// conferences, confError := confrepo.GetByClient(clientid, 0)

	if sessionErr != nil {
		fmt.Printf("sessions by conferene returns error : %v\n", sessionErr)
		return nil, sessionErr

	}

	vm := viewmodels.SessionListVMRead{}
	vm.ClientID = confDB.ClientID
	vm.ClientName = confDB.ClientName
	vm.ConfID = confDB.ID
	vm.ConfName = confDB.Title
	vm.Sessions = sessions

	return &vm, nil
}

// CreateSession by confid
func (srv *SessionService) CreateSession(SessionData viewmodels.SessionCreateVMWrite) (uuid.UUID, error) {
	// confID, _ := uuid.FromString(ConfData.ID)

	sessionRepo := repositories.Sessions{}
	// conferencerepo := repositories.Conferences{}
	session := models.Session{
		// ID:              confID,
		Title:           SessionData.Title,
		Summary:         SessionData.Summary,
		Details:         SessionData.Details,
		DurationDisplay: SessionData.DurationDisplay,
		IsActive:        SessionData.IsActive,
		StartDate:       SessionData.StartDate,
		ClientID:        SessionData.ClientID,
		ConferenceID:    SessionData.ConfID,
		EndDate:         SessionData.EndDate,
		Address:         SessionData.Address,
		GeoLocation:     models.GeoLocation{GeoLocationLat: SessionData.Latitude, GeoLocationLong: SessionData.Longitude, Radius: SessionData.LocationRadius},
		Venue:           SessionData.Venue,
		SortOrder:       SessionData.SortOrder,
		IsFeatured:      SessionData.IsFeatured,
	}
	// conference.ID = confID

	sessionError := sessionRepo.Create(&session)

	if sessionError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", sessionError)
		return uuid.Nil, sessionError

	}
	return session.ID, nil
}

// UpdateSession by confid
func (srv *SessionService) UpdateSession(SessionData viewmodels.SessionCreateVMWrite) (uuid.UUID, error) {
	// confID, _ := uuid.FromString(ConfData.ID)

	fmt.Printf("about to update session in side service. \n")
	fmt.Printf("session passed = %v \n", SessionData)

	sessionRepo := repositories.Sessions{}
	// conferencerepo := repositories.Conferences{}
	SessionModel, ERR := sessionRepo.GetByID(SessionData.ID)
	if ERR != nil {
		fmt.Printf("session by id : %v\n", ERR)
		return uuid.Nil, ERR

	}
	// session := models.Session{
	// 	// ID:              confID,
	SessionModel.Title = SessionData.Title
	SessionModel.Summary = SessionData.Summary
	SessionModel.Details = SessionData.Details
	SessionModel.DurationDisplay = SessionData.DurationDisplay
	SessionModel.IsActive = SessionData.IsActive
	SessionModel.StartDate = SessionData.StartDate
	SessionModel.ClientID = SessionData.ClientID
	SessionModel.ConferenceID = SessionData.ConfID
	SessionModel.EndDate = SessionData.EndDate
	SessionModel.Address = SessionData.Address
	SessionModel.IsFeatured = SessionData.IsFeatured
	SessionModel.Venue = SessionData.Venue
	SessionModel.SortOrder = SessionData.SortOrder
	SessionModel.GeoLocation = models.GeoLocation{GeoLocationLat: SessionData.Latitude, GeoLocationLong: SessionData.Longitude, Radius: SessionData.LocationRadius}
	SessionModel.ID = SessionData.ID
	// conference.ID = confID

	sessionError := sessionRepo.Update(SessionModel)

	if sessionError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", sessionError)
		return uuid.Nil, sessionError

	}
	return SessionModel.ID, nil
}

//GetConferenceAndClient will return viewmodel
func (srv *SessionService) GetConferenceAndClient(ConferenceID string) (*viewmodels.SessionCreateVMRead, error) {
	fmt.Printf("conferenceid  passed in session GetConferenceAndClient is = %v \n", ConferenceID)

	confid, _ := uuid.FromString(ConferenceID)

	// clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	// sessionRepo := repositories.Sessions{}
	//confService := ConferenceService{}
	clientRepo := repositories.Clients{}
	// conferencerepo := repositories.Conferences{}

	// clientDB, err := clientrepo.GetByID(confid)

	// if err != nil {
	// 	fmt.Printf("cant find clientdb =  %v \n", err)
	// 	return nil, err

	// }

	confDB, confError := confrepo.GetByID(confid)
	if confError != nil {
		fmt.Printf("conferece cant be found by id =  %v \n", confError)
	}
	clientDB, clientError := clientRepo.GetByID(confDB.ClientID)
	if clientError != nil {
		fmt.Printf("conferece cant be found by id =  %v \n", clientError)
	}

	// sessions, sessionErr := sessionRepo.GetByConference(confid, 0)

	// // conferences, confError := confrepo.GetByClient(clientid, 0)

	// if sessionErr != nil {
	// 	fmt.Printf("sessions by conferene returns error : %v\n", sessionErr)
	// 	return nil, sessionErr

	// }

	fmt.Println("hye oye clientDB",clientDB.Name)

	vm := viewmodels.SessionCreateVMRead{}
	vm.ClientID = confDB.ClientID
	vm.ClientName = clientDB.Name
	vm.ConfID = confDB.ID
	vm.ConfName = confDB.Title

	return &vm, nil
}

// get session list by speaker id

func (srv *SessionService) GetSessionsBySpeakerID(speakerID string, ConferenceID string) (*viewmodels.SessionListVMRead, error) {
	fmt.Printf("confid passed is = %v \n", ConferenceID)
	fmt.Printf("speakerid passed is = %v \n", speakerID)

	confid, _ := uuid.FromString(ConferenceID)
	spkrid, _ := uuid.FromString(speakerID)

	sessionRepo := repositories.Sessions{}
	confService := ConferenceService{}
	usersrepo := repositories.Users{}

	confDB, confError := confService.GetConferenceByID(ConferenceID)
	if confError != nil {
		fmt.Printf("conferece cant be found by id =  %v \n", confid)
	}

	sessions, sessionErr := sessionRepo.GetBySpeakerId(spkrid, confid)

	if sessionErr != nil {
		fmt.Printf("sessions by speaker returns error : %v\n", sessionErr)
		return nil, sessionErr

	}

	Speaker, sessionErr := usersrepo.GetByID(spkrid)
	if sessionErr != nil {
		fmt.Printf("sessions by speaker returns error : %v\n", sessionErr)
		return nil, sessionErr

	}

	vm := viewmodels.SessionListVMRead{}
	vm.ClientID = confDB.ClientID
	vm.ClientName = confDB.ClientName
	vm.ConfID = confDB.ID
	vm.ConfName = confDB.Title
	vm.Sessions = sessions
	vm.SpeakerFirstName = Speaker.FirstName
	vm.SpeakerLastName = Speaker.LastName
	vm.SpeakerID = Speaker.ID

	return &vm, nil
}
