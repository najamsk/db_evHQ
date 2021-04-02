package repositories

import (

	// Import GORM-related packages.

	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"

	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	uuid "github.com/satori/go.uuid"
)

// Conferences will deal with client model.
type Sponsors struct {
	database.DataBaseManager
	// DB *gorm.DB
}

func (repo *Sponsors) GetByConferenceID(confID uuid.UUID) ([]viewmodels.SponsorVm, error) {
	fmt.Printf("confID passed in sponsor GetByConferenceID = %v \n", confID)
	db := repo.GetDB()

	fmt.Println(&db)
	var sponsor []viewmodels.SponsorVm
	fmt.Printf("confID passed in sponsor GetByConferenceID = %v \n", confID)

	var query string = `select sponsors.id, sponsors.name,sponsors.is_active, sponsors.sort_order,sl.name as type
	from sponsors 
	inner join sponsor_levels sl on sponsors.sponsor_level_id = sl.id 
	 where sponsors.conference_id=? order by sponsors.sort_order;`

	dbErr := db.Raw(query, confID).Find(&sponsor).Error
	if dbErr != nil {
		fmt.Printf("Error at GetByConferenceID repo  in sponsor, = %v \n", dbErr)
		return nil, dbErr
	}
	fmt.Println("hye 9991", sponsor)
	return sponsor, nil
}
func (repo *Sponsors) GetByID(sponID uuid.UUID) (*models.Sponsor, error) {
	fmt.Printf("sponsor id passed is = %v \n", sponID)
	db := repo.GetDB()

	fmt.Println(&db)
	sponsor := models.Sponsor{}
	fmt.Printf("sponsor id passed is= %v \n", sponID)

	dbErr := db.Where("id=?", sponID).Find(&sponsor).Error
	if dbErr != nil {
		fmt.Printf("Error at GetByID in sponsor, = %v \n", dbErr)
		return nil, dbErr
	}
	fmt.Println("hye 9991", sponsor)
	return &sponsor, nil
}
func (repo *Sponsors) UpdateSponsor(UserObj *models.Sponsor) error {
	// Print out the balances.
	db := repo.GetDB()
	defer db.Close()
	errSave := db.Save(UserObj).Error

	if errSave != nil {
		fmt.Println("edr error hye", errSave)
		return errSave
	}

	return nil

}
func (repo *Sponsors) GetSponsorLevel(sponLevelID uuid.UUID) (*models.SponsorLevel, error) {
	fmt.Printf("sponsor level id passed is = %v \n", sponLevelID)
	db := repo.GetDB()

	fmt.Println(&db)
	sponsor := models.SponsorLevel{}
	fmt.Printf("sponsor level id passed is  = %v \n", sponLevelID)

	dbErr := db.Where("id=?", sponLevelID).Find(&sponsor).Error
	if dbErr != nil {
		fmt.Printf("GetSponsorLevel repo  in sponsor, = %v \n", dbErr)
		return nil, dbErr
	}
	fmt.Println("hye 9991", sponsor)
	return &sponsor, nil
}
func (repo *Sponsors) GetALLSponsorLevel() ([]models.SponsorLevel,error) {
	var result []models.SponsorLevel
	db := repo.GetDB()
	defer db.Close()
	ERr:=db.Find(&result).Error
	if ERr !=nil{
		fmt.Println("GetALLSponsorLevel shows error")
		return nil,ERr
	}
	return result,nil
}

func (repo *Sponsors) Create(UserObj *models.Sponsor)(uuid.UUID,error ){
	// Print out the balances.
	db := repo.GetDB()
	defer db.Close()
	errSave := db.Save(UserObj).Error

	if errSave != nil {
		fmt.Println("edr error hye", errSave)
		return uuid.Nil,errSave
	}

	return UserObj.ID,nil

}