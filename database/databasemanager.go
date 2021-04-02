package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/revel/revel"
)

// DataBaseManager wraps around database specs
type DataBaseManager struct {
}

//GetDB would return gormDB with defer close before and last line will be reutrn will it work
func (base *DataBaseManager) GetDB() *gorm.DB {
	addr := revel.Config.StringDefault("hq.database.connectionstring", "")
	if addr == "" {
		fmt.Printf("datbase manager has no valid connection string \n")
	}

	fmt.Printf("openning db = %v\n", addr)
	db, err := gorm.Open("postgres", addr)
	//db.LogMode(true)
	if err != nil {
		fmt.Println("error from db returned")
		fmt.Printf("%v", err)
		log.Fatal(err)
	}

	// defer db.Close()
	return db
}
