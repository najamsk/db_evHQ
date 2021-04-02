package services

import (
	"fmt"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	//"github.com/revel/revel"
)

// AccountService will do dirty work
type AccountService struct{}

//GetByEmail return user by email
func (srv *AccountService) GetByEmail(email string) (*models.User, error) {
	fmt.Printf("GetConferenceByID func got clientid passed  = %v \n", email)

	// clientrepo := repositories.Clients{}
	// confrepo := repositories.Conferences{}
	userrepo := repositories.Users{}
	// conferencerepo := repositories.Conferences{}

	user, err := userrepo.GetByEmail(email)

	if err != nil {
		fmt.Printf("error locating user from account service =  %v \n", err)

		return nil, err

	}

	return user, nil
	// return nil, nil
}

//GetByEmail return user by email
func (srv *AccountService) HQRoles() (map[string]struct{}, error) {
	// fmt.Printf("GetConferenceByID func got clientid passed  = %v \n", email)

	// clientrepo := repositories.Clients{}
	// confrepo := repositories.Conferences{}
	userrepo := repositories.Users{}
	// conferencerepo := repositories.Conferences{}

	hqroles, err := userrepo.GetHQRoles()

	if err != nil {
		fmt.Printf("accountserivce/cant load hq roles from repo =  %v \n", err)
		return nil, err
	}

	return hqroles, nil
	// return nil, nil
}
func (srv *AccountService) UpdateResetPassword(Email string,code string)  error{
	userRepo:=repositories.Accounts{}
	var passCode models.ResetPassword
	passCode.Email=Email
	passCode.Code=code

	RstErr:=userRepo.Update(passCode)
	if RstErr!=nil{
		fmt.Println("UpdateResetPassword repo shows error")
		return RstErr
	}
	
	return nil
}
