package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open("postgres", "user=postgres dbname=todos password=123qwe123 sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&Todo{}, &User{})
	DB = db
}
