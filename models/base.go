package models

import (
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	err := godotenv.Load() //load the .env file
	if err != nil {
		fmt.Println(err)
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println(err)
	}
	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{}) //database migration
}

//Return a handle to the DB object

func GetDB() *gorm.DB {
	return db
}
