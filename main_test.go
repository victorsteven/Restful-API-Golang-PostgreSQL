package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"rest_api_gorm/controllers"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	viper.SetConfigFile("./config.json")
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	testData := viper.Sub("test")
	err := testData.Unmarshal(&conn)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", conn.Host, conn.DbUser, conn.DbName, conn.DbPass)

	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println(err)
	}
	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{}) //database migration
}

type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserId uint   `json:"user_id"` //The user that this contact belongs to
}

type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

var db *gorm.DB

func init() {
	err := godotenv.Load() //load the .env file
	if err != nil {
		fmt.Println(err)
	}
	username := os.Getenv("db_user_test")
	password := os.Getenv("db_pass_test")
	dbName := os.Getenv("db_name_test")
	dbHost := os.Getenv("db_host_test")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println(err)
	}
	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{}) //database migration
}

func clearTable() {
	db.Exec("DELETE FROM accounts")
	db.Exec("ALTER SEQUENCE accounts_id_seq RESTART WITH 1")
}

// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	a.Router.ServeHTTP(rr, req)
// 	return rr
// }

// func checkResponseCode(t *testing.T, expected, actual int) {
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
// 	}
// }

func TestGetHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetHomePage)
	handler.ServeHTTP(rr, req)
	//Check the status code if is what we want:
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := map[string]interface{}{
		"message": "Welcome to the home page", "status": true,
	}
	expectedBinary, err := json.Marshal(expected)
	if err != nil {
		log.Fatal("Could not convert the map to json")
	}
	expectedJSON := string(expectedBinary)
	// Remove the newline character as the last element in the response string
	responseString := strings.TrimSuffix(rr.Body.String(), "\n")
	assert.Equal(t, responseString, expectedJSON)
}

func TestCreateAccount(t *testing.T) {
	// clearTable()

	payload := []byte(`{"email": "gobimermaar@gmail.com", "password": "password"}`)
	req, err := http.NewRequest("POST", "/api/user/new", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("This is the error %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateAccount)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := map[string]interface{}{
		"account": map[string]interface{}{
			"email": "gobimermaar@gmail.com",
		},
		"message": "Account has been created",
		"status":  true,
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &jsonMap)
	if err != nil {
		log.Fatal("Could not convert the map to json")
	}
	assert.Equal(t, jsonMap["message"], expected["message"])
	assert.Equal(t, jsonMap["status"], expected["status"])
	assert.Equal(t, jsonMap["account"].(map[string]interface{})["email"], expected["account"].(map[string]interface{})["email"])

	// Check for duplicate emails
	// payload2 := []byte(`{"email": "gobimermaar@gmail.com", "password": "password"}`)
	// req2, err2 := http.NewRequest("POST", "/api/user/new", bytes.NewBuffer(payload2))
	// if err2 != nil {
	// 	log.Fatalf("This is the error %v", err2)
	// }
	// rr2 := httptest.NewRecorder()
	// handler2 := http.HandlerFunc(controllers.CreateAccount)
	// handler2.ServeHTTP(rr2, req2)

	// if status2 := rr2.Code; status2 != http.StatusOK {
	// 	t.Errorf("Handler returned wrong status code: got %v want %v", status2, http.StatusOK)
	// }
	// expected2 := map[string]interface{}{
	// 	"account": map[string]interface{}{
	// 		"email": "gobimermaar@gmail.com",
	// 	},
	// 	"message": "Email address already in use",
	// 	"status":  false,
	// }
	// jsonMap2 := make(map[string]interface{})
	// err2 = json.Unmarshal([]byte(rr2.Body.String()), &jsonMap2)
	// if err2 != nil {
	// 	log.Fatal("Could not convert the map to json")
	// }
	// assert.Equal(t, jsonMap2["message"], expected2["message"])
	// assert.Equal(t, jsonMap2["status"], expected2["status"])
	// assert.Equal(t, jsonMap["account"].(map[string]interface{})["email"], expected["account"].(map[string]interface{})["email"])
}
