package main

import (
	"fmt"
	"log"

	"rest_api_gorm/controllers"

	"github.com/spf13/viper"
)

type Connection struct {
	Host   string
	DbUser string
	DbName string
	DbPass string
}

var conn Connection

func main() {

	a := controllers.App{}

	viper.SetConfigFile("./config.json")
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	developmentData := viper.Sub("development")
	err := developmentData.Unmarshal(&conn)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	a.Initialize(conn.Host, conn.DbUser, conn.DbPass, conn.DbName)

	a.Run(":8000")
}
