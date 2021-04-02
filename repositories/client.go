package repositories

import (

	// Import GORM-related packages.

	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

// Clients will deal with client model.
type Clients struct {
	database.DataBaseManager

	// DB *gorm.DB
}

// GetAll return all users from repo
func (repo *Clients) GetAll(Skip int, Limit ...int) []models.Client {

	// Print out the balances.
	var clients []models.Client
	db := repo.GetDB()
	defer db.Close()
	if len(Limit) > 0 {
		qLimit := Limit[0]
		db.Offset(Skip).Limit(qLimit).Find(&clients)
	} else {
		db.Offset(Skip).Find(&clients)
	}
	// fmt.Printf("sf %v \n", len(Limit))

	// for _, client := range clients {
	// 	fmt.Printf("%s\n", client.Name)
	// }
	return clients
}

// GetByID return all users from repo
func (repo *Clients) GetByID(ClientID uuid.UUID) (*models.Client, error) {

	// Print out the balances.
	fmt.Println("here in client and client id", ClientID)
	clientDB := models.Client{}
	db := repo.GetDB()
	defer db.Close()
	err := db.Where("id = ?", ClientID).Find(&clientDB).Error

	if err != nil {
		log.Fatal("cant find client by ids1")
		return nil, err
	}

	return &clientDB, nil
}

// Update return all users from repo
func (repo *Clients) Update(ClientID uuid.UUID, Name string, IsActive bool) error {
	clientDB := models.Client{}
	db := repo.GetDB()
	defer db.Close()
	// get conference and then update

	// errDB := db.Where("id=?", ClientID).Find(&clientDB).Error
	// if errDB != nil {
	// 	fmt.Printf("can't find Client by id, error = %v \n", errDB)
	// 	if gorm.IsRecordNotFoundError(errDB) {
	// 		errDB = db.Save(clientDB).Error // newUser not user
	// 	}
	// 	return errDB
	// }
	// Print out the balances.

	clientDB.Name = Name
	clientDB.ID = ClientID
	clientDB.IsActive = IsActive

	err := db.Save(clientDB).Error

	return err
}

// Create return all users from repo
func (repo *Clients) Create(ClientObj models.Client) (uuid.UUID, error) {

	// Print out the balances.
	// var clients []models.Client
	db := repo.GetDB()
	defer db.Close()
	err := db.Create(&ClientObj).Error

	if err != nil {
		log.Fatal(err)
		// fmt.Errorf(err.Error())
		fmt.Printf("create func in client repos throws error = %v \n", err)
		return uuid.Nil, err
	}

	// for _, client := range clients {
	// 	fmt.Printf("%s\n", client.Name)
	// }
	return ClientObj.ID, nil
}
func (repo *Clients) SearchByName(name string, Skip int, Limit ...int) ([]models.Client, error) {

	// Print out the balances.
	var clients []models.Client
	db := repo.GetDB()
	defer db.Close()

	var query string = `select * from clients 
			 where LOWER(clients.name) like ?;`

	dberr := db.Raw(query, "%"+name+"%").Scan(&clients).Error
	if dberr != nil {
		fmt.Println("sreach client shows error", dberr)
		return clients, dberr
	}
	fmt.Println("yeh clients", clients)
	return clients, nil
}
