package controllers

import (
	"encoding/json"
	"net/http"
	"rest_api_gorm/models"
	u "rest_api_gorm/utils"
)

//CreateAccount function
func (a *App) CreateAccount(res http.ResponseWriter, req *http.Request) {
	account := &models.Account{} //instantiate a account object
	err := json.NewDecoder(req.Body).Decode(account)
	if err != nil {
		u.Respond(res, u.Message(false, "Error while decoding the body"))
		return
	}
	resp := account.Create(a.DB)
	u.Respond(res, resp)
}

//Authenticate function
func (a *App) UserLogin(res http.ResponseWriter, req *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(req.Body).Decode(account)
	if err != nil {
		u.Respond(res, u.Message(false, "Invalid request"))
		return
	}
	resp := models.Login(a.DB, account.Email, account.Password)
	u.Respond(res, resp)
}
