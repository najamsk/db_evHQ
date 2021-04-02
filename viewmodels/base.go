package viewmodels

import (
	_ "fmt"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Base contains common columns for all tables.
type Base struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	//ID        int `gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"update_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	CreatedBy *uuid.UUID `gorm:"type:uuid;"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid;"`
	DeletedBy *uuid.UUID `gorm:"type:uuid;"`

	// CreatedBy	int `gorm:"primary_key"`
	// UpdatedBy	int `gorm:"primary_key"`
	// DeletedBy	int `gorm:"primary_key"`
	Deleted bool
}

/* BeforeCreate will set a UUID rather than numeric ID.
If you want to user auto generated uuid as primary key
then uncomment following code */

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid)
}
