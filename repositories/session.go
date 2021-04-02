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

// Sessions will deal with client model.
type Sessions struct {
	database.DataBaseManager
	// DB *gorm.DB
}

// GetAll return all users from repo
func (repo *Sessions) GetAll(Skip int, Limit ...int) []models.Session {

	// Print out the balances.
	var result []models.Session
	db := repo.GetDB()
	defer db.Close()
	if len(Limit) > 0 {
		qLimit := Limit[0]
		db.Offset(Skip).Limit(qLimit).Find(&result)
	} else {
		db.Offset(Skip).Find(&result)
	}
	// fmt.Printf("sf %v \n", len(Limit))

	// for _, client := range clients {
	// 	fmt.Printf("%s\n", client.Name)
	// }
	return result
}

// GetByConference return all users from repo
func (repo *Sessions) GetByConference(ConferenceID uuid.UUID, Skip int, Limit ...int) ([]viewmodels.SessionListRead, error) {

	// Print out the balances.
	var result []viewmodels.SessionListRead

	db := repo.GetDB()
	defer db.Close()
	var query string=`select id,title,start_date,end_date,is_active,sort_order,(select count(*) FROM  session_speakers where session_id=sessions.id) as speaker_count from sessions where sessions.conference_id=? ORDER BY sessions.sort_order DESC`
	if len(Limit) > 0 {
		qLimit := Limit[0]
		dbErr:=db.Raw(query,ConferenceID).Offset(Skip).Limit(qLimit).Scan(&result).Error
		if dbErr !=nil{
			fmt.Println("sesion shows error",dbErr)
		}		//db.Where("conference_id = ?", ConferenceID).Offset(Skip).Limit(qLimit).Find(&result)
	} else {
		dbErr:=db.Raw(query,ConferenceID).Offset(Skip).Scan(&result).Error
		
		if dbErr !=nil{
			fmt.Println("sesion shows error",dbErr)
		}
		//db.Where("conference_id = ?", ConferenceID).Offset(Skip).Find(&result)
	}
	// fmt.Printf("sf %v \n", len(Limit))

	fmt.Printf("sessions return by conf from sessionRepo = %v\n", result)

	return result, nil
}

//Update take conf object and update in db
func (repo *Sessions) Update(Obj *models.Session) error {
	// Print out the balances.
	db := repo.GetDB()
	defer db.Close()
	errSave := db.Save(Obj).Error
	if errSave != nil {
		return errSave
	}

	return nil

}

//Create take conf object and update in db
func (repo *Sessions) Create(Obj *models.Session) error {
	// Print out the balances.
	db := repo.GetDB()
	defer db.Close()

	errSave := db.Create(Obj).Error

	if errSave != nil {
		fmt.Println("Error at create session", errSave)
		return errSave
	}

	return nil

}

func (repo *Sessions) GetByID(SessionID uuid.UUID) (*models.Session, error) {

	// Print out the balances.
	result := models.Session{}
	db := repo.GetDB()
	defer db.Close()

	errDB := db.Where("id=?", SessionID).Find(&result).Error
	if errDB != nil {
		fmt.Printf("can't find session by id, error = %v \n", errDB)
		return nil, errDB
	}

	fmt.Printf("subscription located in model = %v\n", result)

	return &result, nil
}
func (repo *Sessions) GetBySpeakerId(SpeakerID uuid.UUID, ConferenceId uuid.UUID) ([]viewmodels.SessionListRead, error) {
	db := repo.GetDB()
	var result []viewmodels.SessionListRead
	var query string = `select sessions.id,sessions.title,sessions.start_date,sessions.end_date,sessions.is_active,sessions.sort_order,
	(select count(*) FROM  session_speakers where session_id=sessions.id) as speaker_count from sessions  inner join session_speakers cs
	 on sessions.id = cs.session_id   where sessions.conference_id=? and cs.user_id=? ORDER BY sessions.sort_order DESC`
	sessionErr := db.Raw(query, ConferenceId, SpeakerID).Scan(&result).Error
	if sessionErr != nil {
		fmt.Println("Error at session by speaker id ", sessionErr)
		return nil, sessionErr
	}
	return result, nil
}
