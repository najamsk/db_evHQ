package services

import (
	"fmt"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	uuid "github.com/satori/go.uuid"
)

type SponsorService struct {
}

func (srv *SponsorService) GetByconference(confid string) (*viewmodels.SponsorVmRead, error) {
	conferenceID, _ := uuid.FromString(confid)
	spnRepo := repositories.Sponsors{}
	confRepo := repositories.Conferences{}
	clientRepo := repositories.Clients{}

	sponsors, ERR := spnRepo.GetByConferenceID(conferenceID)
	if ERR != nil {
		fmt.Println("sponser by conference shows err", ERR)
		return nil, ERR

	}
	conference, confERR := confRepo.GetByID(conferenceID)
	if ERR != nil {
		fmt.Println("confRepo.GetByID shows err", confERR)
		return nil, ERR

	}
	ClientDB, ERR := clientRepo.GetByID(conference.ClientID)
	if ERR != nil {
		fmt.Println("clientRepo.GetByID shows err", ERR)
		return nil, ERR

	}
	sponserLeveDB, LeveDBERR := spnRepo.GetALLSponsorLevel()
	if ERR != nil {
		fmt.Println("spnRepo.GetALLSponsorLevel shows err", LeveDBERR)
		return nil, LeveDBERR

	}
	vm := viewmodels.SponsorVmRead{}
	vm.ClientID = ClientDB.ID
	vm.ClientName = ClientDB.Name
	vm.Sponsors = sponsors
	vm.ConferenceID = conferenceID
	vm.Sponsorlevel = sponserLeveDB

	return &vm, nil
}

func (srv *SponsorService) GetByID(id string, confid string) (*viewmodels.SponsorEditVmRead, error) {
	sponsorID, _ := uuid.FromString(id)
	confID, _ := uuid.FromString(confid)
	spnRepo := repositories.Sponsors{}
	confRepo := repositories.Conferences{}
	clientRepo := repositories.Clients{}

	sponsors, ERR := spnRepo.GetByID(sponsorID)
	if ERR != nil {
		fmt.Println("sponser by ID shows err", ERR)
		return nil, ERR

	}
	sponsorLvl, SpnErr := spnRepo.GetSponsorLevel(sponsors.SponsorLevelID)
	if ERR != nil {
		fmt.Println("sponserleve by ID shows err", SpnErr)
		return nil, ERR

	}
	conference, confERR := confRepo.GetByID(confID)
	if ERR != nil {
		fmt.Println("sponser by conference shows err", confERR)
		return nil, ERR

	}
	ClientDB, ERR := clientRepo.GetByID(conference.ClientID)
	if ERR != nil {
		fmt.Println("sponser by conference shows err", ERR)
		return nil, ERR

	}
	sponserLeveDB, LeveDBERR := spnRepo.GetALLSponsorLevel()
	if ERR != nil {
		fmt.Println("sponser by conference shows err", LeveDBERR)
		return nil, LeveDBERR

	}
	vm := viewmodels.SponsorEditVmRead{}
	vm.ClientID = ClientDB.ID
	vm.ClientName = ClientDB.Name
	vm.Name = sponsors.Name
	vm.Description = sponsors.Description
	vm.IsActive = sponsors.IsActive
	vm.SortOrder=sponsors.SortOrder
	vm.ID = sponsors.ID
	vm.ConferenceID = confID
	vm.Facebook = sponsors.SocialMedia.Facebook
	vm.Youtube = sponsors.SocialMedia.Youtube
	vm.Twitter = sponsors.SocialMedia.Twitter
	vm.Linkedin = sponsors.SocialMedia.LinkedIn
	vm.Sponsorlevel_ID = sponsorLvl.ID
	vm.SponsorLevel_name = sponsorLvl.Name
	vm.Sponsorlevel = sponserLeveDB

	return &vm, nil
}
func (srv *SponsorService) Update(spnData viewmodels.SponsorEditVmWrite) (uuid.UUID, error) {
	//clientID, _ := uuid.FromString("8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e")
	// speakerID, _ := uuid.FromString(speakerData.ID)
	spnRepo := repositories.Sponsors{}
	spnModel, usrError := spnRepo.GetByID(spnData.ID)
	if usrError != nil {
		fmt.Printf("spnRepo.GetByID returns error : %v\n", usrError)
		return uuid.Nil, usrError

	}
	spnModel.Name = spnData.Name
	spnModel.Description = spnData.Bio
	spnModel.ClientID = spnData.ClientID
	spnModel.IsActive = spnData.IsActive
	spnModel.ID = spnData.ID
	spnModel.SponsorLevelID = spnData.Sponlevel
	spnModel.SocialMedia.Facebook = spnData.Facebook
	spnModel.SocialMedia.Twitter = spnData.Twitter
	spnModel.SocialMedia.LinkedIn = spnData.Linkedin
	spnModel.SocialMedia.Youtube = spnData.Youtube
	spnModel.SortOrder = spnData.SortOrder
	spnError := spnRepo.UpdateSponsor(spnModel)
	if spnError != nil {
		fmt.Printf("spnRepo.UpdateSponsor returns error : %v\n", spnError)
		return uuid.Nil, spnError

	}
	fmt.Printf("updating sponsor id :", spnData.ID)

	return spnData.ID, nil
}
func (srv *SponsorService) Create(spnData viewmodels.SponsorEditVmWrite) (uuid.UUID, error) {
	spnRepo := repositories.Sponsors{}
	spnModel := models.Sponsor{
		Name:           spnData.Name,
		Description:    spnData.Bio,
		ClientID:       spnData.ClientID,
		IsActive:       spnData.IsActive,
		SponsorLevelID: spnData.Sponlevel,
		SortOrder:      spnData.SortOrder,
		ConferenceID:   spnData.ConferenceID,
	}
	spnModel.SocialMedia.Facebook = spnData.Facebook
	spnModel.SocialMedia.Twitter = spnData.Twitter
	spnModel.SocialMedia.LinkedIn = spnData.Linkedin
	spnModel.SocialMedia.Youtube = spnData.Youtube
	spnID, spnError := spnRepo.Create(&spnModel)
	if spnError != nil {
		fmt.Printf("spnRepo.Create returns error : %v\n", spnError)
		return uuid.Nil, spnError

	}
	fmt.Printf("updating sponsor id :", spnData.ID)

	return spnID, nil
}
