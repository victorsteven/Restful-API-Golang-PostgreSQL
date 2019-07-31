package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"rest_api_gorm/controllers"
	"testing"
)

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

	//Check the response body is what we expect.
	// expected := struct{
	// 	"message": "Welcome to the home page",
	// 	"status":  true,
	// }
	expected := map[string]interface{}{
		"message": "Welcome to the home page",
		"status":  true,
	}

	if rr.Body.String() != expected["message"] {
		t.Errorf("handler returned enexpected body: got %v want %v", rr.Body.String(), expected)
	}
	// log.Println("this is the body: ", rr.Body[message])
}
