package controllers

import (
	"encoding/json"
	"net/http"
	"rest_api_gorm/models"
	u "rest_api_gorm/utils"
)

func CreateAccount(res http.ResponseWriter, req *http.Request) {

	account := &models.Account{} //instantiate a account object

	// fmt.Printf("This is the first result: %v\n", req.Body)

	err := json.NewDecoder(req.Body).Decode(account)
	if err != nil {
		u.Respond(res, u.Message(false, "Error while decoding the body"))
		return
	}
	resp := account.Create()
	u.Respond(res, resp)
}

// var Getme = func(res http.ResponseWriter, req *http.Request) {
// 	resp := map[string]string{"hello bro": "Yes"}
// 	u.Respond(res, resp)
// }

func Authenticate(res http.ResponseWriter, req *http.Request) {
	account := &models.Account{}
	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(req.Body).Decode(account)
	if err != nil {
		u.Respond(res, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(res, resp)
}
