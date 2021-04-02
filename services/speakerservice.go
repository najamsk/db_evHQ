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
type SpeakerService struct {
}

func (srv *SpeakerService) GetSpeakerByConferenceID(conferenceID string) (*viewmodels.SpeakerListVMRead, error) {
	fmt.Printf("conferenceID in speaker GetSpeakerByConferenceID is = %v \n", conferenceID)

	confid, _ := uuid.FromString(conferenceID)

	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	speakerrepo := repositories.Speaker{}
	// conferencerepo := repositories.Conferences{}

	conferenceDB, confError := confrepo.GetByID(confid)
	if confError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", confError)
		return nil, confError

	}
	clientDB, err := clientrepo.GetByID(conferenceDB.ClientID)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}
	speakers, err := speakerrepo.GetByConferenceID(confid)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	vm := viewmodels.SpeakerListVMRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.ConferenceID = conferenceDB.ID
	vm.Speakers = speakers

	return &vm, nil
}
func (srv *SpeakerService) GetSpeakerByID(speakerID string, sessionID string) (*viewmodels.SpeakerEditVMRead, error) {
	fmt.Printf("GetConferenceByID func got clientid passed  = %v \n", speakerID)

	Speakerid, _ := uuid.FromString(speakerID)
	sessionid, _ := uuid.FromString(sessionID)
	sessionrepo := repositories.Sessions{}
	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	speakerrepo := repositories.Speaker{}
	imgrepo := repositories.Images{}

	speaker, spkerr := speakerrepo.GetByID(Speakerid)

	if spkerr != nil {
		fmt.Printf("cant find clientdb =  %v \n", spkerr)
		return nil, spkerr

	}
	SessionDB, seesionerr := sessionrepo.GetByID(sessionid)
	if seesionerr != nil {
		fmt.Printf("cant find sessiondb =  %v \n", seesionerr)
		return nil, seesionerr
	}

	ConfDB, conferr := confrepo.GetByID(SessionDB.ConferenceID)
	if conferr != nil {
		fmt.Printf("cant find confdb =  %v \n", conferr)
		return nil, conferr
	}

	fmt.Println("here in speaker service clientid", ConfDB.ClientID, "heye", speaker)
	clientDB, err := clientrepo.GetByID(ConfDB.ClientID)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	sessionSpkr, sessionErr := speakerrepo.SessionSpeakerByid(Speakerid, sessionid)
	if sessionErr != nil {
		fmt.Printf("cant find clientdb =  %v \n", sessionErr)
		return nil, sessionErr

	}
	confSpkr, confErr := speakerrepo.ConferenceSpeakerByid(Speakerid, SessionDB.ConferenceID)
	if confErr != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, confErr

	}
	vm := viewmodels.SpeakerEditVMRead{}
	Profile, prerr := imgrepo.GetImage(speaker.ID, "user", "user_profile")
	if prerr == nil {
		vm.ProfileURL = Profile.BasicURL + Profile.ImageURLPrefix + "/" + Profile.Name

	}
	poster, poserr := imgrepo.GetImage(speaker.ID, "user", "poster")
	if poserr == nil {
		vm.PosterURL = poster.BasicURL + poster.ImageURLPrefix + "/" + poster.Name

	}

	vm.ClientID = clientDB.ID
	vm.ConfID = SessionDB.ConferenceID
	vm.ClientName = clientDB.Name
	vm.ID = speaker.ID
	vm.FirstName = speaker.FirstName
	vm.LastName = speaker.LastName
	vm.Email = speaker.Email
	vm.Designation = speaker.Designation
	vm.PhoneNumber = speaker.PhoneNumber
	vm.Bio = speaker.Bio
	vm.SessionID = sessionid
	vm.Facebook = speaker.SocialMedia.Facebook
	vm.Twitter = speaker.SocialMedia.Twitter
	vm.Linkedin = speaker.SocialMedia.LinkedIn
	vm.Youtube = speaker.SocialMedia.Youtube
	vm.Organization = speaker.Organization
	vm.ConferenceWeight = confSpkr.SortOrder
	vm.SessionWeight = sessionSpkr.SortOrder

	return &vm, nil
	// return nil, nil
}
func (srv *SpeakerService) UpdateSpeaker(speakerData viewmodels.SpeakerEditVMWrite) (uuid.UUID, error) {
	//clientID, _ := uuid.FromString("8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e")
	// speakerID, _ := uuid.FromString(speakerData.ID)
	usersRepo := repositories.Users{}
	speakerRepo := repositories.Speaker{}
	SpeakerModel, ERR := speakerRepo.GetByID(speakerData.ID)
	if ERR != nil {
		fmt.Printf("speaker by id returns error : %v\n", ERR)
		return uuid.Nil, ERR

	}
	// conferencerepo := repositories.Conferences{}

	SpeakerModel.FirstName = speakerData.FirstName
	SpeakerModel.LastName = speakerData.LastName
	SpeakerModel.Email = speakerData.Email
	SpeakerModel.Organization = speakerData.Organization
	SpeakerModel.Designation = speakerData.Designation
	SpeakerModel.PhoneNumber = speakerData.PhoneNumber
	SpeakerModel.Bio = speakerData.Bio
	SpeakerModel.ClientID = speakerData.ClientID
	SpeakerModel.ID = speakerData.ID
	SpeakerModel.SocialMedia.Facebook = speakerData.Facebook
	SpeakerModel.SocialMedia.Twitter = speakerData.Twitter
	SpeakerModel.SocialMedia.LinkedIn = speakerData.Linkedin
	SpeakerModel.SocialMedia.Youtube = speakerData.Youtube
	usrID, UsrError := usersRepo.Update(SpeakerModel)
	fmt.Printf("updating speaker user id :", usrID)

	if UsrError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", UsrError)
		return uuid.Nil, UsrError

	}
	sessionSpeaker := models.Session_speakers{
		SessionID: speakerData.SessionID,
		UserID:    usrID,
		SortOrder: speakerData.SessionWeight,
	}
	sessionError := speakerRepo.AddSessionSpeaker(&sessionSpeaker)
	if sessionError != nil {
		fmt.Println("error in usr repo", UsrError)
		return uuid.Nil, sessionError
	}

	conferenceSpeaker := models.Conference_speakers{
		ConferenceID: speakerData.ConferenceID,
		UserID:       usrID,
		SortOrder:    speakerData.ConferenceWeight,
	}
	confError := speakerRepo.AddConferenceSpeaker(&conferenceSpeaker)
	if confError != nil {
		fmt.Println("error in usr repo", confError)
		return uuid.Nil, confError
	}
	return usrID, nil
}
func (srv *SpeakerService) GetSpeakerBySession(id string, confid string) (*viewmodels.SpeakerListVMRead, error) {
	fmt.Printf("confid passed is = %v \n", confid)

	confiD, _ := uuid.FromString(confid)
	sessionID, _ := uuid.FromString(id)

	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	speakerrepo := repositories.Speaker{}
	sessionrepo := repositories.Sessions{}

	conferenceDB, confError := confrepo.GetByID(confiD)
	if confError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", confError)
		return nil, confError

	}
	clientDB, err := clientrepo.GetByID(conferenceDB.ClientID)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}
	speakers, err := speakerrepo.GetBySessionID(sessionID)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	session, err := sessionrepo.GetByID(sessionID)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	vm := viewmodels.SpeakerListVMRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.ConferenceID = conferenceDB.ID
	vm.Speakers = speakers
	vm.SessionID = sessionID
	vm.SessionTitle = session.Title

	return &vm, nil
}
func (srv *SpeakerService) AddsessionSpeaker(obj viewmodels.AddSessionSpeakerVMWrite) (uuid.UUID, error) {
	fmt.Println("hyeeee user", obj)
	usersRepo := repositories.Users{}
	speakerRepo := repositories.Speaker{}
	var SpeakerModel *models.User
	var usrID uuid.UUID
	var ERR error
	var UsrError error

	// conferencerepo := repositories.Conferences{}
	if obj.UserId != "" {
		fmt.Println("userid", obj.UserId)
		UserID, _ := uuid.FromString(obj.UserId)
		// user.ID = UserID
		SpeakerModel, ERR = speakerRepo.GetByID(UserID)
		if ERR != nil {
			fmt.Println("Speaker not found by id", ERR)
			return uuid.Nil, ERR
		}
		SpeakerModel.FirstName = obj.FirstName
		SpeakerModel.LastName = obj.LastName
		SpeakerModel.Email = obj.Email
		SpeakerModel.Organization = obj.Organization
		SpeakerModel.Designation = obj.Designation
		SpeakerModel.PhoneNumber = obj.PhoneNumber
		SpeakerModel.Bio = obj.Bio
		SpeakerModel.SocialMedia.Facebook = obj.Facebook
		SpeakerModel.SocialMedia.Twitter = obj.Twitter
		SpeakerModel.SocialMedia.LinkedIn = obj.Linkedin
		SpeakerModel.SocialMedia.Youtube = obj.Youtube
		SpeakerModel.IsActive = true
		usrID, UsrError = usersRepo.Update(SpeakerModel)
		if UsrError != nil {
			fmt.Println("error in usr repo", UsrError)
			return uuid.Nil, UsrError
		}

		fmt.Println("useriduuuid", UserID)
	} else {
		Speaker := models.User{
			FirstName:    obj.FirstName,
			LastName:     obj.LastName,
			Email:        obj.Email,
			Organization: obj.Organization,
			Designation:  obj.Designation,
			PhoneNumber:  obj.PhoneNumber,
			Bio:          obj.Bio,
			IsActive:     true,
		}
		Speaker.SocialMedia.Facebook = obj.Facebook
		Speaker.SocialMedia.Twitter = obj.Twitter
		Speaker.SocialMedia.LinkedIn = obj.Linkedin
		Speaker.SocialMedia.Youtube = obj.Youtube

		usrID, UsrError = usersRepo.Update(&Speaker)
		if UsrError != nil {
			fmt.Println("error in usr repo", UsrError)
			return uuid.Nil, UsrError
		}
	}
	sessionSpeaker := models.Session_speakers{
		SessionID: obj.SessionID,
		UserID:    usrID,
		SortOrder: obj.SessionWeight,
	}
	sessionError := speakerRepo.AddSessionSpeaker(&sessionSpeaker)
	if sessionError != nil {
		fmt.Println("error in usr repo", sessionError)
		return uuid.Nil, sessionError
	}

	conferenceSpeaker := models.Conference_speakers{
		ConferenceID: obj.ConfID,
		UserID:       usrID,
		SortOrder:    obj.ConferenceWeight,
	}
	confError := speakerRepo.AddConferenceSpeaker(&conferenceSpeaker)
	if confError != nil {
		fmt.Println("error in usr repo", confError)
		return uuid.Nil, confError
	}

	return usrID, nil
}
