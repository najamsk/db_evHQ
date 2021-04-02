package repositories

import (

	// Import GORM-related packages.

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"
	//"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
)

// Conferences will deal with client model.
type Accounts struct {
	database.DataBaseManager
	// DB *gorm.DB
}

func (repo *Accounts) Update(passwordCode models.ResetPassword) error {
	db := repo.GetDB()
	fmt.Println("passwordCode1:", passwordCode)
	var codedb models.ResetPassword
	err := db.Where("email = ?", passwordCode.Email).First(&codedb).Error

	if err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			fmt.Println("insert")
			fmt.Println("passwordCode2:", passwordCode)
			err = db.Create(&passwordCode).Error // newUser not user
		}

	} else {
		fmt.Println("update")
		err = db.Model(passwordCode).Update(&passwordCode).Error
		fmt.Println("err:", err)
	}

	return err
}

func (repo *Accounts) GetPasswordTokenExpiryInMin(email string, code string) (int,error) {
	db := repo.GetDB()

	type Result struct {
		ExpiryMin int
	}

	var result Result

	var query string = `SELECT cast(cast((now() - updated_at) as int)/60 as int) AS expiry_min
						FROM "reset_passwords" where email = ? and code = ?`

	err := db.Raw(query, email, code).Scan(&result).Error
	fmt.Println("result:", result)
	if err != nil {
		return -1 ,err
	}

	return result.ExpiryMin,nil
}
