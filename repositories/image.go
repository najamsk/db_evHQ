package repositories

import (

	// Import GORM-related packages.

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

// Sessions will deal with client model.
type Images struct {
	database.DataBaseManager
	// DB *gorm.DB
}

func (repo *Images) Create(Obj *models.Image) error {
	// Print out the balances.
	fmt.Println("##############################################here in create image", Obj)
	db := repo.GetDB()
	defer db.Close()

	errSave := db.Create(Obj).Error

	if errSave != nil {
		return errSave
	}

	return nil

}
func (repo *Images) GetByEntityid(entityID uuid.UUID) []models.Image {
	var images []models.Image
	db := repo.GetDB()
	defer db.Close()
	db.Where("EntityId = ?", entityID).Find(&images)
	fmt.Printf("confernces located in model = %v\n", images)

	return images
}
func (repo *Images) GetImage(entityID uuid.UUID, entityType string, imageCategory string) (models.Image, error) {
	db := repo.GetDB()
	var imageobj models.Image
	var err = db.Where("entity_id = ? AND entity_type = ? AND image_category = ?", entityID, entityType, imageCategory).First(&imageobj).Error
	if err != nil {
		return imageobj, err
	}

	return imageobj, nil
}
func (repo *Images) Update(imageobj models.Image) (models.Image, error) {
	db := repo.GetDB()
	result := models.Image{}

	var err = db.Where("entity_id = ? AND entity_type = ? AND image_category = ?", imageobj.EntityId, imageobj.EntityType, imageobj.ImageCategory).Find(&result).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			errCreate := db.Save(&imageobj).Error // newUser not user
			if errCreate != nil {
				return imageobj, errCreate
				fmt.Println("create image shows err", errCreate)
			}
			fmt.Println("image successfully created")
			return imageobj, nil
		}
		fmt.Println("find image shows err", err)
		return imageobj, err

	}

	imageobj.ID = result.ID

	fmt.Println("hye imGE ID", imageobj.ID)

	errSave := db.Save(imageobj).Error
	if errSave != nil {
		return imageobj, errSave
	}
	fmt.Println("image successfully updated")

	return imageobj, nil
}
