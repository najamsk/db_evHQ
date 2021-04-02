package models

import(
	 "github.com/jinzhu/gorm"
	 "github.com/satori/go.uuid"
	"path/filepath")

// User type
type User struct {
	Base
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Organization string
	Designation  string
	ProfileImg   string
	PhoneNumber	 string
	PhoneNumber2 string
	Bio			 string
	SocialMedia  SocialMedia `gorm:"embedded"`
	IsActive     bool
	GeoLocation  GeoLocation `gorm:"embedded"`
	ClientID     uuid.UUID
	Conferences  []*Conference `gorm:"many2many:users_conferences;"`
	Roles        []*Role       `gorm:"many2many:users_roles;"`
	MyAgenda     MyAgenda
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
		
	uuid, err := uuid.NewV4()
	if err != nil {
	 return err
	}

	//if profile image is not empty then set/replace ProfileImg column value with user primary key which is uuid 
	if(user.ProfileImg != ""){
		scope.SetColumn("ProfileImg", uuid.String()+ filepath.Ext(user.ProfileImg))
	}
	return scope.SetColumn("ID", uuid)
   } 

