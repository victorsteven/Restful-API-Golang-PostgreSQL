package controllers

import (
	"fmt"
	"log"
	"net/http"
	"rest_api_gorm/app"
	"rest_api_gorm/models"
	u "rest_api_gorm/utils"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App struct
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(host, user, password, dbname string) {
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
	var err error
	a.DB, err = gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println(err)
	}
	a.DB.Debug().AutoMigrate(&models.Account{}, &models.Contact{}) //database migration

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	fmt.Println("Listening to port 8000")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", GetHomePage).Methods("GET")
	a.Router.HandleFunc("/api/user/new", a.CreateAccount).Methods("POST")
	a.Router.HandleFunc("/api/user/login", a.UserLogin).Methods("POST")
	a.Router.HandleFunc("/api/contacts/new", a.CreateContact).Methods("POST")
	a.Router.HandleFunc("/api/{id}/contacts", a.GetContactsFor).Methods("GET")

	a.Router.Use(app.JwtAuthentication) //attach JWT auth middleware
}

//GetHomePage function
func GetHomePage(res http.ResponseWriter, req *http.Request) {
	u.Respond(res, u.Message(true, "Welcome to the home page"))
}
