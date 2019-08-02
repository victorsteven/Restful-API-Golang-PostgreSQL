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

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// type Contact struct {
// 	gorm.Model
// 	Name   string `json:"name"`
// 	Phone  string `json:"phone"`
// 	UserId uint   `json:"user_id"` //The user that this contact belongs to
// }

// type Account struct {
// 	gorm.Model
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// 	Token    string `json:"token";sql:"-"`
// }

var a controllers.App

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
	a = controllers.App{}

	fmt.Printf("this is the database name: %v\n", conn.DbName)

	a.Initialize(conn.Host, conn.DbUser, conn.DbPass, conn.DbName)

	// ensureTableExists()

	code := m.Run()

	os.Exit(code)
}

func clearTable() {
	a.DB.Exec("DELETE FROM accounts")
	a.DB.Exec("ALTER SEQUENCE accounts_id_seq RESTART WITH 1")
}

// func TestEmptyTable(t *testing.T) {
// 	clearTable()
// 	req, _ := http.NewRequest("GET", "/users", nil)
// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// 	if body := response.Body.String(); body != "[]" {
// 		t.Errorf("Expected an empty array. Got %s", body)
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
	clearTable()

	payload := []byte(`{"email": "pet@gmail.com", "password": "password"}`)
	req, err := http.NewRequest("POST", "/api/user/new", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("This is the error %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(a.CreateAccount)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := map[string]interface{}{
		"account": map[string]interface{}{
			"email": "pet@gmail.com",
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

	//Check for duplicate emails
	payload2 := []byte(`{"email": "pet@gmail.com", "password": "password"}`)
	req2, err2 := http.NewRequest("POST", "/api/user/new", bytes.NewBuffer(payload2))
	if err2 != nil {
		log.Fatalf("This is the error %v", err2)
	}
	rr2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(a.CreateAccount)
	handler2.ServeHTTP(rr2, req2)

	if status2 := rr2.Code; status2 != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status2, http.StatusOK)
	}
	expected2 := map[string]interface{}{
		"account": map[string]interface{}{
			"email": "pet@gmail.com",
		},
		"message": "Email address already in use",
		"status":  false,
	}
	jsonMap2 := make(map[string]interface{})
	err2 = json.Unmarshal([]byte(rr2.Body.String()), &jsonMap2)
	if err2 != nil {
		log.Fatal("Could not convert the map to json")
	}
	assert.Equal(t, jsonMap2["message"], expected2["message"])
	assert.Equal(t, jsonMap2["status"], expected2["status"])
	assert.Equal(t, jsonMap["account"].(map[string]interface{})["email"], expected["account"].(map[string]interface{})["email"])
}

func TestUserLogin(t *testing.T) {
	clearTable()
	payload := []byte(`{"email": "pet@gmail.com", "password": "password"}`)
	req, err := http.NewRequest("POST", "/api/user/new", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("This is the error %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(a.CreateAccount)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := map[string]interface{}{
		"account": map[string]interface{}{
			"email": "pet@gmail.com",
		},
		"message": "Account has been created",
		"status":  true,
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &jsonMap)
	if err != nil {
		log.Fatal("Could not convert the map to json")
	}
	userToken := jsonMap["account"].(map[string]interface{})["token"]
	userTokenString := fmt.Sprintf("%v", userToken) //Converting the interface token to string
	assert.Equal(t, jsonMap["message"], expected["message"])
	assert.Equal(t, jsonMap["status"], expected["status"])
	assert.Equal(t, jsonMap["account"].(map[string]interface{})["email"], expected["account"].(map[string]interface{})["email"])

	loginPayload := []byte(`{"email": "pet@gmail.com", "password": "password"}`)

	reqLogin, errLogin := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(loginPayload))
	if errLogin != nil {
		log.Fatalf("This is the error %v", errLogin)
	}
	// Set the authentication token here:
	reqLogin.Header.Set("Authorization", "Bearer "+userTokenString)

	rrLogin := httptest.NewRecorder()
	handlerLogin := http.HandlerFunc(a.UserLogin)
	handlerLogin.ServeHTTP(rrLogin, reqLogin)

	if statusLogin := rrLogin.Code; statusLogin != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", statusLogin, http.StatusOK)
	}
	expectedLogin := map[string]interface{}{
		"account": map[string]interface{}{
			"email": "pet@gmail.com",
		},
		"message": "Logged In",
		"status":  true,
	}
	jsonMapLogin := make(map[string]interface{})
	errLogin = json.Unmarshal([]byte(rrLogin.Body.String()), &jsonMapLogin)
	if errLogin != nil {
		log.Fatal("Could not convert the map to json")
	}
	assert.Equal(t, jsonMapLogin["message"], expectedLogin["message"])
	assert.Equal(t, jsonMapLogin["status"], expectedLogin["status"])
	assert.Equal(t, jsonMapLogin["account"].(map[string]interface{})["email"], expectedLogin["account"].(map[string]interface{})["email"])
}
