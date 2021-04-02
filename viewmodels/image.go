
package viewmodels

import (
	//"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

type ImageVMWrite struct {
	Base
	Name                 string
	BasicURL             string
	FolderPath           string
	ImageURLPrefix       string 
	EntityID             uuid.UUID
	EntityType           string
	ImageCategory        string
	IsActive             bool
}