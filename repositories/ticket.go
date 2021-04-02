package repositories

import (

	// Import GORM-related packages.

	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
)

type Tickets struct {
	database.DataBaseManager
	// DB *gorm.DB
}

func (repo *Tickets) CreateTicket(Obj *models.Ticket) error {
	db := repo.GetDB()
	defer db.Close()
	TicketModel := models.Ticket{}
	fmt.Println("ticket obj passed in create ticket", Obj)
	errDB := db.Where("serial_no=?", Obj.SerialNo).Find(&TicketModel).Error
	if errDB != nil {
		if gorm.IsRecordNotFoundError(errDB) {
			fmt.Println("no err")
			errDb := db.Save(Obj).Error // newUser not user
			if errDB != nil {
				fmt.Println("yahan error hye createticket repo main", errDb)
				return errDb
			} else {
				return nil
			}
			return errDb

		}
		fmt.Println("errors at CreateType repo", errDB)
		return errDB
	}
	fmt.Println("Seriol number already exist", errDB)
	return errors.New("Already exist")
}

func (repo *Tickets) GetByConference(ConferenceID uuid.UUID) ([]models.Ticket, error) {

	// Print out the balances.
	var result []models.Ticket

	db := repo.GetDB()
	defer db.Close()
	DbErr := db.Where("conference_id = ?", ConferenceID).Order("valid_from ASC").Find(&result).Error
	if DbErr != nil {
		fmt.Println("Error in GetByConference repositry in tickets", DbErr)
		return nil, DbErr
	}
	// fmt.Printf("sf %v \n", len(Limit))

	fmt.Printf("sessions return by conf from sessionRepo = %v\n", result)

	return result, nil
}
func (repo *Tickets) GetById(TicketId uuid.UUID) (*models.Ticket, error) {
	Ticket := models.Ticket{}
	db := repo.GetDB()
	defer db.Close()
	errDB := db.Where("id=?", TicketId).Find(&Ticket).Error
	if errDB != nil {
		fmt.Printf("can't find tickets by id, error = %v \n", errDB)
		return nil, errDB
	}
	return &Ticket, nil
}
func (repo *Tickets) UpdateTicket(Obj *models.Ticket) error {
	db := repo.GetDB()
	defer db.Close()

	errSave := db.Save(Obj).Error
	if errSave != nil {
		fmt.Printf("ticket update, error = %v \n", errSave)
		return errSave
	}
	return nil
}
