package models

// _ "github.com/jinzhu/gorm/dialects/postgres"
// _ "github.com/go-sql-driver/mysql"

// type App struct {
// 	Router *mux.Router
// 	DB     *gorm.DB
// }

// func (a *App) Initialize(host, user, password, dbname string) {
// 	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
// 	var err error
// 	a.DB, err = gorm.Open("postgres", dbUri)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	a.DB.Debug().AutoMigrate(&Account{}, &Contact{}) //database migration

// 	a.Router = mux.NewRouter()
// 	a.initializeRoutes()

// }

// func (a *App) initializeRoutes() {
// 	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
// 	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
// 	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
// 	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
// 	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
// }

// type Connection struct {
// 	Host   string
// 	DbUser string
// 	DbName string
// 	DbPass string
// }

// var conn Connection

// func init() {
// 	viper.SetConfigFile("./config.json")
// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Fatalf("Error reading config file, %s", err)
// 	}
// 	fmt.Printf("Using config dev: %s\n", viper.ConfigFileUsed())
// 	developmentData := viper.Sub("development")
// 	err := developmentData.Unmarshal(&conn)
// 	if err != nil {
// 		log.Fatalf("unable to decode into struct, %v", err)
// 	}
// 	// InitConnection(conn.Host, conn.DbUser, conn.DbName, conn.DbPass)
// }
