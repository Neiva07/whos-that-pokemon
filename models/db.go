package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //drive to import mysql dialect
	"github.com/joho/godotenv"
)

//Database - struct holding the db itself and the methods attached to it
type Database struct {
	db *gorm.DB
}

//Initialize method to initialize the database
func (DB *Database) Initialize() {

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	// dbHost := os.Getenv("db_host")

	dbURI := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbName)

	connection, err := gorm.Open("mysql", dbURI)

	if err != nil {
		log.Print(err)
	}
	DB.db = connection

	DB.db.Debug().AutoMigrate(&User{}, &Generation{}, &Game{}, &Friendship{})

}

//GetDB exporting the dabaase to the rest of the application
func (DB *Database) GetDB() *gorm.DB {
	return DB.db
}

//DB declare the variable that could be exported
var DB Database

func init() {

	err := godotenv.Load(os.ExpandEnv("$GOPATH/src/whos-that-pokemon/.env"))

	if err != nil {
		log.Print(err)
	}
	DB.Initialize()
	log.Println("DB Initialized!")

}
