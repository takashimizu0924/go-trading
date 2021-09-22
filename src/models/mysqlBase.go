package models

import (
	"app/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var AppDB *sql.DB

const (
	tableNameSignalEvents = "signal_events"
)

func GetCandleTableName(product_code string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", product_code, duration)
}

func NewMysqlBase() {
	db, err := gorm.Open(mysql.Open(config.Config.Mysql), &gorm.Config{})
	if err != nil {
		log.Println("databaseOpenError", err)
	}
	cmd := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (time DATETIME PRIMARY KEY NOT NULL,product_code varchar(10),side varchar(10),price float,size float)`, tableNameSignalEvents)
	db.Exec(cmd)
	for _, duration := range config.Config.Durations {
		tableName := GetCandleTableName("BTC_JPY", duration)
		c := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (time DATETIME PRIMARY KEY NOT NULL,open float,close float,hight float,low float,volume float)`, tableName)
		db.Exec(c)
	}
	// db.Table("buy_orders").Migrator().CreateTable(&OrderEvent{})

	// db.Table("sell_orders").Migrator().CreateTable(&OrderEvent{})
	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfull got Mysql DB connection!!!")
	AppDB = sqlDB
}
