package repositories

import (

	// Import GORM-related packages.

	"fmt"
	//"log"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	uuid "github.com/satori/go.uuid"
)

type TicketsType struct {
	database.DataBaseManager
}

func (repo *TicketsType) GetByConference(ConferenceID uuid.UUID) ([]viewmodels.TicketTypelist, error) {

	// Print out the balances.
	var result []viewmodels.TicketTypelist

	 db := repo.GetDB()
	defer db.Close()
	// dbErr:=db.Where("conference_id = ?", ConferenceID).Find(&result).Error
	// if dbErr != nil {
	// 	fmt.Printf("Error at GetByConferenceID repo , = %v \n", dbErr)
	// 	return nil, dbErr
	// }


	var query string = `select ticket_types.id,ticket_types.title,ticket_types.Client_id,ticket_types.is_active,ticket_types.conference_id,
	ticket_types.amount, ticket_types.ammount_currency ,ticket_types.description,
	(select count(*) FROM  tickets where ticket_type_id=ticket_types.id) as total_tickect,
	(select count(*) FROM  tickets where ticket_type_id=ticket_types.id and is_consumed='true' ) as consumed_ticket from ticket_types
	  where ticket_types.conference_id=?`
	//errDB := db.Where("id=?", TicketTypeId).Find(&TicketType).Error
	dbErr := db.Raw(query,ConferenceID).Scan(&result).Error
	
	if dbErr != nil {
		fmt.Printf("Error at GetByConferenceID repo , = %v \n", dbErr)
		return nil, dbErr
	}
	
	// fmt.Printf("sf %v \n", len(Limit))

	fmt.Printf("sessions return by conf from sessionRepo = %v\n", result)

	return result, nil
}
func (repo *TicketsType) CreateTicketType(Obj *models.TicketType) error {
	db := repo.GetDB()
	defer db.Close()

	errSave := db.Create(Obj).Error

	if errSave != nil {
		fmt.Println("cant create ticket type in ticket repo")
		return errSave
	}

	return nil
}
func (repo *TicketsType) GetTypeById(TicketTypeId uuid.UUID) (*models.TicketType, error) {
	TicketType := models.TicketType{}
	db := repo.GetDB()
	defer db.Close()
	errDB := db.Where("id=?", TicketTypeId).Find(&TicketType).Error
	if errDB != nil {
		fmt.Printf("can't find tickettype by id, error = %v \n", errDB)
		return nil, errDB
	}
	return &TicketType, nil
}
func (repo *TicketsType) UpdateTiketType(Obj *models.TicketType)error{
	db :=repo.GetDB()
	defer db.Close()
	errSave := db.Save(Obj).Error
	if errSave != nil {
		fmt.Printf("tickettype update, error = %v \n", errSave)
			return errSave
	}
	return nil
}
