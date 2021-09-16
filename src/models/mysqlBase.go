package models

import (
	"app/config"
	"database/sql"
	"log"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var AppDB *sql.DB

func NewMysqlBase() {
	db, err := gorm.Open(mysql.Open(config.Config.Mysql), &gorm.Config{})
	if err != nil {
		log.Println("databaseOpenError", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfull got Mysql DB connection!!!")
	AppDB = sqlDB
}
