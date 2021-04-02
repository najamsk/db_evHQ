package models

import(
	// "github.com/jinzhu/gorm"
	 "github.com/satori/go.uuid"
	//"path/filepath"
)

type Users_roles struct {
	Base
	UserID            uuid.UUID
	RoleID	  		  uuid.UUID

}