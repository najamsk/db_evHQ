package repositories

import (

	// Import GORM-related packages.

	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

// Conferences will deal with client model.
type Conferences struct {
	database.DataBaseManager
	// DB *gorm.DB
}

// GetAll return all users from repo
func (repo *Conferences) GetAll(Skip int, Limit ...int) []models.Subscription {

	// Print out the balances.
	var subscriptions []models.Subscription
	db := repo.GetDB()
	defer db.Close()
	if len(Limit) > 0 {
		qLimit := Limit[0]
		db.Offset(Skip).Limit(qLimit).Find(&subscriptions)
	} else {
		db.Offset(Skip).Find(&subscriptions)
	}
	// fmt.Printf("sf %v \n", len(Limit))

	// for _, client := range clients {
	// 	fmt.Printf("%s\n", client.Name)
	// }
	return subscriptions
}

// GetByClient return all users from repo
func (repo *Conferences) GetByClient(ClientID uuid.UUID, Skip int, Limit ...int) ([]models.Conference, error) {

	// Print out the balances.
	var confernces []models.Conference

	db := repo.GetDB()
	defer db.Close()
	if len(Limit) > 0 {
		qLimit := Limit[0]
		db.Where("client_id = ?", ClientID).Offset(Skip).Limit(qLimit).Find(&confernces)
	} else {
		db.Where("client_id = ?", ClientID).Offset(Skip).Find(&confernces)
	}
	// fmt.Printf("sf %v \n", len(Limit))

	fmt.Printf("confernces located in model = %v\n", confernces)

	return confernces, nil
}

//UpdateConference take conf object and update in db
func (repo *Conferences) UpdateConference(ConferenceObj *models.Conference) error {
	// Print out the balances.
	db := repo.GetDB()
	defer db.Close()

	errSave := db.Save(ConferenceObj).Error

	if errSave != nil {
		return errSave
	}

	return nil

}

//CreateConference take conf object and update in db
func (repo *Conferences) CreateConference(ConferenceObj *models.Conference) error {
	// Print out the balances.
	db := repo.GetDB()
	defer db.Close()

	errSave := db.Create(ConferenceObj).Error

	if errSave != nil {
		fmt.Println("Errors at CreateConference",errSave)
		return errSave
	}

	return nil

}

// GetByClient return all users from repo
func (repo *Conferences) GetByID(ConferenceID uuid.UUID) (models.Conference, error) {

	// Print out the balances.
	conf := models.Conference{}
	db := repo.GetDB()
	defer db.Close()
	fmt.Println("ConferenceID",ConferenceID)
	errDB := db.Where("id=?", ConferenceID).Find(&conf).Error
	if errDB != nil {
		fmt.Printf("can't find conference by id, error = %v \n", errDB)
		return conf, errDB
	}

	fmt.Printf("subscription located in model = %v\n", conf)

	return conf, nil
}

// // Create return all users from repo
// func (repo *Conferences) Create(SubscriptionObj *models.Subscription) (uuid.UUID, error) {

// 	// Print out the balances.
// 	// var clients []models.Client
// 	db := repo.GetDB()
// 	defer db.Close()
// 	err := db.Create(SubscriptionObj).Error

// 	if err != nil {
// 		log.Fatal(err)
// 		// fmt.Errorf(err.Error())
// 		return uuid.Nil, err
// 	}

// 	// for _, client := range clients {
// 	// 	fmt.Printf("%s\n", client.Name)
// 	// }
// 	return SubscriptionObj.ID, nil
// }
