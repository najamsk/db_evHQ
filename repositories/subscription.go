package repositories

import (

	// Import GORM-related packages.

	"errors"
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

// Clients will deal with client model.
type Subscriptions struct {
	database.DataBaseManager
	// DB *gorm.DB
}

// GetAll return all users from repo
func (repo *Subscriptions) GetAll(Skip int, Limit ...int) []models.Subscription {

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

// GetAll return all users from repo
func (repo *Subscriptions) GetSubscriptionByClient(ClientID uuid.UUID, Skip int, Limit ...int) (*models.Subscription, error) {

	// Print out the balances.
	var subscription models.Subscription
	db := repo.GetDB()
	defer db.Close()
	if len(Limit) > 0 {
		qLimit := Limit[0]
		db.Where("client_id = ?", ClientID).Offset(Skip).Limit(qLimit).Preload("Payments").Find(&subscription)
	} else {
		db.Where("client_id = ?", ClientID).Offset(Skip).Preload("Payments").Find(&subscription)
	}
	// fmt.Printf("sf %v \n", len(Limit))

	fmt.Printf("subscription located in model = %v\n", subscription)
	fmt.Printf("sub id %T \n", subscription.ID)
	fmt.Printf("sub id %v \n", subscription.ID)

	zerouuid, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

	subequal := uuid.Equal(zerouuid, subscription.ID)
	if subequal == true {
		fmt.Println("Errors at GetSubscriptionByClient repo ",subequal)
		return nil, errors.New("no subscription found")
	}
	fmt.Printf("subequal check = %v \n", subequal)
	// for _, client := range clients {
	// 	fmt.Printf("%s\n", client.Name)
	// }
	return &subscription, nil
}

// GetByID return all users from repo
func (repo *Subscriptions) GetByID(SubscriptionID uuid.UUID) (*models.Subscription, error) {

	// Print out the balances.
	subscriptionDB := models.Subscription{}
	db := repo.GetDB()
	defer db.Close()
	err := db.Where("id = ?", SubscriptionID).Preload("Payments").Find(&subscriptionDB).Error

	if err != nil {
		log.Fatal("cant find client by id")
		return nil, err
	}

	return &subscriptionDB, nil
}

func (repo *Subscriptions) UpdateSubscription(SubscriptionObj *models.Subscription) bool {

	// Print out the balances.
	// var clients []models.Client
	db := repo.GetDB()
	defer db.Close()
	err := db.Save(SubscriptionObj).Error

	if err != nil {
		fmt.Printf("error updating subscription= %v \n", err)
		return false
	}

	// for _, client := range clients {
	// 	fmt.Printf("%s\n", client.Name)
	// }
	return true
}

// Create return all users from repo
func (repo *Subscriptions) Create(SubscriptionObj *models.Subscription) (uuid.UUID, error) {

	// Print out the balances.
	// var clients []models.Client
	db := repo.GetDB()
	defer db.Close()
	err := db.Create(SubscriptionObj).Error

	if err != nil {
		log.Fatal(err)
		// fmt.Errorf(err.Error())
		return uuid.Nil, err
	}

	// for _, client := range clients {
	// 	fmt.Printf("%s\n", client.Name)
	// }
	return SubscriptionObj.ID, nil
}
