package services

import (
	"fmt"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"

	//"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"

	//"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// ConferenceService will do dirty work
type UserService struct {
}

func (srv *UserService) Getlargestrole(role []models.Role) models.Role {
	userrole := models.Role{}
	if len(role) < 1 {
		return userrole
	}
	max := role[0].Weight // assume first value is the smallest
	userrole = role[0]    //largest role
	for _, value := range role {
		if value.Weight > max {
			max = value.Weight // found another smaller value, replace previous value in max
			userrole = value
		}
	}

	return userrole
}
func (srv *UserService) GetByID(userid string, logeinID uuid.UUID) (*viewmodels.UserEditVMRead, error) {
	roleSrv := RoleService{}

	fmt.Printf("GetById func got userid passed  = %v \n", userid)

	userID, err := uuid.FromString(userid)
	if err != nil {
		fmt.Println("cant convert uuid of user", err)
		return nil, err
	}
	Rolerepo := repositories.Roles{}
	userrepo := repositories.Users{}
	logedinRoles, err := Rolerepo.GetByuserID(logeinID)
	if err != nil {
		fmt.Println("cant get logedin roles", err)
		return nil, err

	}

	userDB, err := userrepo.GetByID(userID)
	if err != nil {
		fmt.Println("cant get user by id", err)
		return nil, err

	}
	userRoles, err := Rolerepo.GetByuserID(userID)
	if err != nil {
		fmt.Println("cant get logedin roles", err)
		return nil, err

	}
	// largest role weight of user
	userLargestRole := roleSrv.Getlargestrole(userRoles)
	// largest role weight of loginuser
	loginLargetRole := roleSrv.Getlargestrole(logedinRoles)

	vm := viewmodels.UserEditVMRead{}
	vm.ID = userDB.ID
	vm.FirstName = userDB.FirstName
	vm.LastName = userDB.LastName
	vm.Email = userDB.Email
	vm.Designation = userDB.Designation
	vm.PhoneNumber = userDB.PhoneNumber
	vm.Bio = userDB.Bio
	vm.Organization = userDB.Organization
	vm.IsActive = userDB.IsActive
	vm.Facebook = userDB.SocialMedia.Facebook
	vm.Twitter = userDB.SocialMedia.Twitter
	vm.Linkedin = userDB.SocialMedia.LinkedIn
	vm.Youtube = userDB.SocialMedia.Youtube
	vm.UserRoleweight = userLargestRole
	vm.LoginRoleweight = loginLargetRole
	return &vm, nil
}
func (srv *UserService) UpdateUser(userData viewmodels.UserEditVMWrite) (uuid.UUID, error) {
	//clientID, _ := uuid.FromString("8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e")
	// speakerID, _ := uuid.FromString(speakerData.ID)
	usersRepo := repositories.Users{}
	roleRepo := repositories.Roles{}
	userModel, usrError := usersRepo.GetByID(userData.ID)
	if usrError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", usrError)
		return uuid.Nil, usrError

	}
	userModel.FirstName = userData.FirstName
	userModel.LastName = userData.LastName
	userModel.Email = userData.Email
	userModel.Organization = userData.Organization
	userModel.Designation = userData.Designation
	userModel.PhoneNumber = userData.PhoneNumber
	userModel.Bio = userData.Bio
	userModel.ClientID = userData.ClientID
	userModel.IsActive = userData.IsActive
	userModel.ID = userData.ID
	userModel.SocialMedia.Facebook = userData.Facebook
	userModel.SocialMedia.Twitter = userData.Twitter
	userModel.SocialMedia.LinkedIn = userData.Linkedin
	userModel.SocialMedia.Youtube = userData.Youtube
	userid, UsrError := usersRepo.Update(userModel)
	if UsrError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", UsrError)
		return uuid.Nil, UsrError

	}
	fmt.Printf("updating user user id :", userid)
	// for i := 0; i < len(userData.Roles); i++ {
	// 	roleid, _ := uuid.FromString(userData.Roles[i])
	// 	fmt.Println("uuid", roleid)
	// 	role.ID = roleid
	//user.Roles = []*models.Role{&role}

	//}
	roleError := roleRepo.UpdateUserRoles(userid, userData.Roles)
	if roleError != nil {
		fmt.Printf("confercnes by client returns error : %v\n", roleError)
		return uuid.Nil, roleError

	}

	return userid, nil
}
