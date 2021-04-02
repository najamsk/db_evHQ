package viewmodels

import uuid "github.com/satori/go.uuid"

//import "github.com/najamsk/eventvisor/eventvisorHQ/models"

type UserEditVMWrite struct {
	//ID string
	Base
	ClientID     uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Organization string
	Designation  string
	PhoneNumber  string
	Bio          string
	IsActive     bool
	Roles        []string
	Facebook         string
	Youtube          string
	Linkedin         string
	Twitter          string
}
type UserRoles struct {
	Name   string
	ID     uuid.UUID
	Isrole bool // just for user role checking
}

type UserEditVMRead struct {
	Base
	ID           uuid.UUID
	FirstName    string
	LastName     string
	PhoneNumber  string
	Bio          string
	Organization string
	Email        string
	Designation  string
	//UserRoles  [] models.Role
	//Loginuserroles[] models.Role
	UserRoleweight  int // largest role weight
	LoginRoleweight int // largest role of login user
	ProfileURL      string
	PosterURL       string
	IsActive        bool
	Roles           []UserRoles // role listing
	Facebook         string
	Youtube          string
	Linkedin         string
	Twitter          string
}
