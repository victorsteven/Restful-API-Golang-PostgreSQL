package models

import (
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
)

type Connection struct {
	Host   string
	DbUser string
	DbName string
	DbPass string
}

var conn Connection

var db *gorm.DB

func init() {
	// a := App{}
	viper.SetConfigFile("./config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	developmentData := viper.Sub("development")
	err := developmentData.Unmarshal(&conn)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	// a.Initialize("steven", "here", "api_medium_kelvin")
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", conn.Host, conn.DbUser, conn.DbName, conn.DbPass)

	// a.Run(":8000")
	// err := godotenv.Load() //load the .env file
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// username := os.Getenv("db_user")
	// password := os.Getenv("db_pass")
	// dbName := os.Getenv("db_name")
	// dbHost := os.Getenv("db_host")

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
