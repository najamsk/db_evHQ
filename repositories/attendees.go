package repositories

import (

	// Import GORM-related packages.

	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"

	//"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	uuid "github.com/satori/go.uuid"
)

// Conferences will deal with client model.
type Attendies struct {
	database.DataBaseManager
	// DB *gorm.DB
}

func (repo *Attendies) GetByConferenceID(confID uuid.UUID) ([]models.User, error) {
	fmt.Printf("confid id passed in attendies = %v \n", confID)
	db := repo.GetDB()

	fmt.Println(&db)
	var attendies []models.User
	fmt.Printf("confid id passed in attendies = %v \n", confID)

	var query string = `select users.id, users.first_name, users.last_name,users.organization,users.designation
	from users 
	inner join conference_attendees ca on users.id = ca.user_id 
	 where ca.conference_id=?;`

	DBErr:=db.Raw(query, confID).Scan(&attendies).Error
	if DBErr!=nil{
		fmt.Println("Errors at GetByConferenceID,",DBErr)
		return nil,DBErr
	}
	fmt.Println("hye 999", attendies)
	return attendies, nil
}
func (repo *Attendies) GetByID(AttendiesID uuid.UUID) (*models.User, error) {
	
	// Print out the balances.
	conf := models.User{}
	db := repo.GetDB()
	defer db.Close()

	errDB := db.Where("id=?", AttendiesID).Find(&conf).Error
	if errDB != nil {
		fmt.Printf("can't find Attendies by id, error = %v \n", errDB)
		return nil, errDB
	}


	return &conf, nil
}