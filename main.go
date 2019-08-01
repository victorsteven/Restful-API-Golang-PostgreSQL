package main

import (
	"fmt"
	"net/http"
	"os"
	"rest_api_gorm/app"
	"rest_api_gorm/controllers"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)

	//Get port from .env file, we did not specify any port so this should return an empty string when tested locally

	router.HandleFunc("/", controllers.GetHomePage).Methods("GET")
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/{id}/contacts", controllers.GetContactsFor).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	fmt.Println("This is the port: ", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Println("This is the error: ", err)
	}
}

func getme(res http.ResponseWriter, req *http.Request) {
	fmt.Println(res, "hello sir")
}
