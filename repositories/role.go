package repositories

import (

	// Import GORM-related packages.

	"fmt"
	"strings"

	//"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisorHQ/database"
	"github.com/najamsk/eventvisor/eventvisorHQ/models"

	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	uuid "github.com/satori/go.uuid"
)

type Roles struct {
	database.DataBaseManager
	// DB *gorm.DB
}

func (repo *Roles) GetByuserID(userid uuid.UUID) ([]models.Role, error) {
	fmt.Printf("userid passed in roles = %v \n", userid)
	db := repo.GetDB()

	fmt.Println(&db)
	var roles []models.Role
	fmt.Printf("userid passed in roles = %v \n", userid)

	var query string = `select roles.id, roles.name, roles.weight
	from roles 
	inner join users_roles ur on roles.id = ur.role_id 
	 where ur.user_id=?;`

	DbErr := db.Raw(query, userid).Scan(&roles).Error
	if DbErr != nil {
		fmt.Println("Errors at GetByuserID repo in role", DbErr)
		return nil, DbErr
	}
	fmt.Println("hye 999", roles)
	return roles, nil
}

func (repo *Roles) GetbyWeight(Id uuid.UUID, weight int) []viewmodels.UserRoles {

	fmt.Printf("weight passed is = %v \n", weight)
	db := repo.GetDB()
	var roles []viewmodels.UserRoles
	var query string
	fmt.Println(&db)
	fmt.Printf("Weight passed in GetbyWeight= %v \n", weight)
	if weight == 1000 {
		query = `select roles.name, roles.id,EXISTS(select user_id FROM users_roles where user_id =? and role_id=roles.id) as Isrole
		from "roles" where weight<=? ORDER BY roles.name`

	} else {
		query = `select roles.name, roles.id,EXISTS(select user_id FROM users_roles where user_id =? and role_id=roles.id) as Isrole
		from "roles" where weight<? ORDER BY roles.name`
	}

	DbErr := db.Raw(query, Id, weight).Scan(&roles).Error
	if DbErr != nil {
		fmt.Println("Errors at GetbyWeight repo in role", DbErr)
	}
	fmt.Println("hye 999", roles)
	return roles
}
func (repo *Roles) UpdateUserRoles(userid uuid.UUID, roleids []string) error {
	db := repo.GetDB()
	fmt.Println("dbConnection2:", &db)
	fmt.Println(&db)

	var query string = `insert into users_roles (user_id, role_id) SELECT ?, ?
						WHERE NOT EXISTS (SELECT * FROM "users_roles" WHERE "user_id" = ? AND "role_id" = ?) `
	//fmt.Println(query)
	for i, v := range roleids {
		fmt.Println(i)
		fmt.Println(v)
		roleid, _ := uuid.FromString(v)

		DbErr := db.Exec(query, userid, roleid, userid, roleid).Error
		if DbErr != nil {
			fmt.Println("error at role update", DbErr)
			return DbErr
		}
	}
	fmt.Printf("shitty roleids: %v\n", roleids)

	ids := strings.Join(roleids, "','")
	deleteQuery := fmt.Sprintf(`delete from users_roles where user_id = ? and role_id not in ('%s')`, ids)

	DbErr := db.Exec(deleteQuery, userid).Error
	if DbErr != nil {
		fmt.Println("error at role deletion", DbErr)
		return DbErr
	}

	return nil
}
