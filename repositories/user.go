package repositories

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Users will deal with user models.
type Users struct {
	database.DataBaseManager
	// DB *gorm.DB
}

//get all users by email
func (repo *Users) GetAllByEmail(Email string) ([]models.User, error) {

	// Print out the balances.
	var result []models.User
	db := repo.GetDB()
	defer db.Close()
	// db.LogMode(true)
	errDB := db.Where("email=?", Email).First(&result).Error
	if errDB != nil {
		fmt.Printf("UserRepo Error: can't find record by id, error = %v \n", errDB)
		return nil, errDB
	}

	fmt.Printf("result by id = %v\n", result)

	return result, nil
}

// GetAll return all users from repo
func (repo *Users) GetAll(Skip int, Limit ...int) []models.User {

	// Print out the balances.
	var result []models.User
	db := repo.GetDB()
	defer db.Close()
	if len(Limit) > 0 {
		qLimit := Limit[0]
		db.Offset(Skip).Limit(qLimit).Find(&result)
	} else {
		db.Offset(Skip).Find(&result)
	}

	return result
}

//Update take conf object and update in db
func (repo *Users) Update(Obj *models.User) (uuid.UUID, error) {

	// Print out the balances.
	db := repo.GetDB()
	defer db.Close()
	//UsrModel.Roles=Obj.Roles

	errSave := db.Save(Obj).Error

	if errSave != nil {
		return uuid.Nil, errSave
	}

	return Obj.ID, nil

}

//Create take conf object and update in db
func (repo *Users) Create(Obj *models.User) error {
	// Print out the balances.
	db := repo.GetDB()
	defer db.Close()

	errSave := db.Create(Obj).Error

	if errSave != nil {
		return errSave
	}

	return nil

}

// GetByID return all users from repo
func (repo *Users) GetByID(userID uuid.UUID) (*models.User, error) {

	// Print out the balances.
	result := models.User{}
	db := repo.GetDB()
	defer db.Close()

	errDB := db.Where("id=?", userID).Find(&result).Error
	if errDB != nil {
		fmt.Printf("can't find record by id, error = %v \n", errDB)
		return nil, errDB
	}

	fmt.Printf("result by id = %v\n", result)

	return &result, nil
}

// GetByEmail return all users from repo
func (repo *Users) GetByEmail(Email string) (*models.User, error) {

	// Print out the balances.
	result := models.User{}
	db := repo.GetDB()
	defer db.Close()

	errDB := db.Where("email=?", Email).Preload("Roles", func(db *gorm.DB) *gorm.DB {
		return db.Order("roles.weight desc")
	}).First(&result).Error
	if errDB != nil {
		fmt.Printf("UserRepo Error: can't find record by id, error = %v \n", errDB)
		return nil, errDB
	}

	fmt.Printf("result by id = %v\n", result)

	return &result, nil
}

// GetHQRoles return all hq level roles
func (repo *Users) GetHQRoles() (map[string]struct{}, error) {

	// Print out the balances.
	result := []models.Role{}
	db := repo.GetDB()
	defer db.Close()

	errDB := db.Where("name like 'HQ%'").Find(&result).Error
	if errDB != nil {
		fmt.Printf("UserRepo Error: can't find record by id, error = %v \n", errDB)
		return nil, errDB
	}

	fmt.Printf("result by id = %v\n", result)
	roleMap := make(map[string]struct{})
	for _, role := range result {
		roleMap[role.Name] = struct{}{}
	}

	return roleMap, nil
}

func (repo *Users) UpdatePassword(email string, password string) error {
	db := repo.GetDB()
	var userObj models.User
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = db.Model(&userObj).Where("email = ?", email).Update("password", string(hash)).Error

	if err != nil {
		return err
	}

	return nil
}
