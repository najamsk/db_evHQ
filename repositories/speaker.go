package repositories

import (

	// Import GORM-related packages.

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"

	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	uuid "github.com/satori/go.uuid"
)

// Conferences will deal with client model.
type Speaker struct {
	database.DataBaseManager
	// DB *gorm.DB
}

func (repo *Speaker) GetByConferenceID(confID uuid.UUID) ([]viewmodels.SpeakerVM, error) {
	fmt.Printf("confid passed in speaker GetByConferenceID is = %v \n", confID)
	db := repo.GetDB()

	fmt.Println(&db)
	var speakers []viewmodels.SpeakerVM
	fmt.Printf("confid passed in speaker GetByConferenceID is = %v \n", confID)

	var query string = `select users.id, users.first_name,users.email, users.last_name,users.organization,users.designation,users.is_active, cs.sort_order
	from users 
	inner join conference_speakers cs on users.id = cs.user_id 
	 where cs.conference_id=? order by cs.sort_order;`

	dbErr := db.Raw(query, confID).Scan(&speakers).Error
	if dbErr != nil {
		fmt.Printf("Error at GetByConferenceID repo , = %v \n", dbErr)
		return nil, dbErr
	}
	fmt.Println("hye 999", speakers)
	return speakers, nil
}

func (repo *Speaker) GetByID(SpeakerID uuid.UUID) (*models.User, error) {

	// Print out the balances.
	conf := models.User{}
	db := repo.GetDB()
	defer db.Close()

	errDB := db.Where("id=?", SpeakerID).Find(&conf).Error
	if errDB != nil {
		fmt.Printf("can't find speaker by id, error = %v \n", errDB)
		return nil, errDB
	}

	fmt.Printf("speaker located in model = %v\n", conf)

	return &conf, nil
}
func (repo *Speaker) UpdateSpeaker(UserObj *models.User) error {
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
func (repo *Speaker) GetBySessionID(id uuid.UUID) ([]viewmodels.SpeakerVM, error) {
	fmt.Printf("sessionid passed in speaker GetBySessionID = %v \n", id)
	db := repo.GetDB()

	fmt.Println(&db)
	var speakers []viewmodels.SpeakerVM
	fmt.Printf("sessionid passed in speaker GetBySessionID= %v \n", id)

	var query string = `select users.id,users.is_active, users.first_name,users.email,spk.sort_order, users.last_name,users.organization,users.designation
	from users 
	inner join session_speakers spk on users.id = spk.user_id 
	 where spk.session_id=? order by spk.sort_order;`

	dbErr := db.Raw(query, id).Scan(&speakers).Error
	if dbErr != nil {
		fmt.Printf("Error at GetByConferenceID repo , = %v \n", dbErr)
		return nil, dbErr
	}
	fmt.Println("hye 999", speakers)
	return speakers, nil
}
func (repo *Speaker) AddSessionSpeaker(spkrObj *models.Session_speakers) error {
	db := repo.GetDB()
	speaker := models.Session_speakers{}
	errDB := db.Where("user_id = ? AND session_id = ?", spkrObj.UserID, spkrObj.SessionID).Find(&speaker).Error
	if errDB != nil {
		fmt.Printf("can't find speaker by id at AddSessionSpeaker, error = %v \n", errDB)
		if gorm.IsRecordNotFoundError(errDB) {
			errDB = db.Save(spkrObj).Error // newUser not user
		}
		return errDB
	}
	speaker.UserID = spkrObj.UserID
	speaker.SessionID = spkrObj.SessionID
	speaker.SortOrder=spkrObj.SortOrder
	DBerr := db.Save(speaker).Error
	if DBerr != nil {
		fmt.Println("cant update session speaker", DBerr)
		return DBerr
	}
	return nil
}

func (repo *Speaker) AddConferenceSpeaker(spkrObj *models.Conference_speakers) error {
	db := repo.GetDB()
	speaker := models.Conference_speakers{}
	errDB := db.Where("user_id = ? AND conference_id = ?", spkrObj.UserID, spkrObj.ConferenceID).Find(&speaker).Error
	if errDB != nil {
		fmt.Printf("can't find speaker by id at AddConferenceSpeaker, error = %v \n", errDB)
		if gorm.IsRecordNotFoundError(errDB) {
			errDB = db.Save(spkrObj).Error // newUser not user
		}
		return errDB
	}
	speaker.UserID = spkrObj.UserID
	speaker.ConferenceID = spkrObj.ConferenceID
	speaker.SortOrder=spkrObj.SortOrder
	DBerr := db.Save(speaker).Error
	if DBerr != nil {
		fmt.Println("cant update conference speaker", DBerr)
		return DBerr
	}
	return nil
}
func(repo *Speaker)SessionSpeakerByid(userID uuid.UUID,SessionID uuid.UUID) (*models.Session_speakers,error){
	conf := models.Session_speakers{}
	db := repo.GetDB()
	defer db.Close()

	errDB := db.Where("user_id = ? AND session_id = ?", userID,SessionID).Find(&conf).Error
	if errDB != nil {
		fmt.Printf("can't find speaker by id, error = %v \n", errDB)
		return nil, errDB
	}

	fmt.Printf("speaker located in model = %v\n", conf)

	return &conf, nil
}
func(repo *Speaker)ConferenceSpeakerByid(userID uuid.UUID,conferenceID uuid.UUID) (*models.Conference_speakers,error){
	conf := models.Conference_speakers{}
	db := repo.GetDB()
	defer db.Close()

	errDB := db.Where("user_id = ? AND conference_id = ?", userID,conferenceID).Find(&conf).Error
	if errDB != nil {
		fmt.Printf("can't find speaker by id, error = %v \n", errDB)
		return nil, errDB
	}

	fmt.Printf("speaker located in model = %v\n", conf)

	return &conf, nil
}
func (repo *Speaker) DeleteSessionSpeaker(SpeakerID uuid.UUID,SessionID uuid.UUID)error{
	db := repo.GetDB()
	defer db.Close()
	deleteQuery := fmt.Sprintf(`delete from session_speakers where user_id = ? and session_id =?`)
	
		DbErr := db.Exec(deleteQuery, SpeakerID,SessionID).Error
		if DbErr != nil {
			fmt.Println("error at role deletion", DbErr)
			return DbErr
		}
return nil
}

