package services

import (
	"fmt"
	//"github.com/najamsk/eventvisor/eventvisorHQ/models"
	//"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"

	//"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// ConferenceService will do dirty work
type AttendiesService struct {
}

func (srv *AttendiesService) GetAttendiesByConferenceID(conferenceID string) (*viewmodels.AttendiesListVMRead, error) {
	fmt.Printf("confid in GetAttendiesByConferenceID = %v \n", conferenceID)

	confid, _ := uuid.FromString(conferenceID)

	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	attendiesrepo := repositories.Attendies{}
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
	attendies, err := attendiesrepo.GetByConferenceID(confid)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	vm := viewmodels.AttendiesListVMRead{}
	vm.ClientID = clientDB.ID
	vm.ClientName = clientDB.Name
	vm.ConferenceID = conferenceDB.ID
	vm.Attendies = attendies

	return &vm, nil
}
func (srv *AttendiesService) GetAttendiesByID(attendiesID string, conferenceiD string) (*viewmodels.AttendiesEditVMRead, error) {
	fmt.Printf("GetConferenceByID func got clientid passed  = %v \n", attendiesID)

	attendiesid, _ := uuid.FromString(attendiesID)
	confiD, _ := uuid.FromString(conferenceiD)

	clientrepo := repositories.Clients{}
	confrepo := repositories.Conferences{}
	attendiesrepo := repositories.Attendies{}
	// conferencerepo := repositories.Conferences{}

	attendie, err := attendiesrepo.GetByID(attendiesid)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}

	confDB, err := confrepo.GetByID(confiD)
	if err != nil {
		fmt.Printf("cant find confdb =  %v \n", err)
		return nil, err
	}
	clientDB, err := clientrepo.GetByID(confDB.ClientID)

	if err != nil {
		fmt.Printf("cant find clientdb =  %v \n", err)
		return nil, err

	}
	vm := viewmodels.AttendiesEditVMRead{}
	vm.ClientID = clientDB.ID
	vm.ConfID = confiD
	vm.ClientName = clientDB.Name
	vm.ID=			attendie.ID
	vm.FirstName=    attendie.FirstName
	vm.LastName=     attendie.LastName
	vm.Email=        attendie.Email
	vm.Designation= attendie.Designation
	vm.PhoneNumber=  attendie.PhoneNumber
	vm.Bio=          attendie.Bio
	vm.Organization= attendie.Organization

	return &vm, nil
	// return nil, nil
}

// UpdateAttendee and return error or nil
func (srv *AttendiesService) UpdateAttendee(AttendeeData viewmodels.AttendiesEditVMWrite) error {
	// confID, _ := uuid.FromString(ConfData.ID)

	fmt.Printf("about to update attendee inside service. \n")
	fmt.Printf("data passed = %v \n", AttendeeData)

	usersRepo := repositories.Users{}

	userModel, usrError := usersRepo.GetByID(AttendeeData.ID)
	if usrError != nil {
		fmt.Printf("Attendie by id returns error : %v\n", usrError)
		return usrError

	}
	userModel.FirstName=    AttendeeData.FirstName
	userModel.LastName =    AttendeeData.LastName
	userModel.Email=        AttendeeData.Email
	userModel.Designation=  AttendeeData.Designation
	userModel.PhoneNumber=  AttendeeData.PhoneNumber
	userModel.Bio=          AttendeeData.Bio
	userModel.Organization= AttendeeData.Organization
	userModel.ID = AttendeeData.ID

	userid,userError := usersRepo.Update(userModel)
	// sessionError := sessionRepo.Update(&session)
	fmt.Printf("updating attendee user id :",userid )

	if userError != nil {
		fmt.Printf("updating attendee from userrepo throws error : %v\n", userError)
		return userError

	}
	return nil
}
