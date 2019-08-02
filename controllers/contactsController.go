package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest_api_gorm/models"
	u "rest_api_gorm/utils"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) CreateContact(res http.ResponseWriter, req *http.Request) {
	//Get the id the user that sents the request
	user_id := req.Context().Value("user_id").(uint)

	//Instantiate a contact object:
	contact := &models.Contact{}
	err := json.NewDecoder(req.Body).Decode(contact)
	if err != nil {
		u.Respond(res, u.Message(false, "Error while decoding request body"))
		return
	}
	contact.UserId = user_id //assign the userId
	resp := contact.Create(a.DB)
	u.Respond(res, resp)
}

func (a *App) GetContactsFor(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	fmt.Println("these are the params: ", params)
	//convert ASCII to Integer
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		//The passed path parameter is not an integer
		u.Respond(res, u.Message(false, "There was an error in your request"))
		return
	}
	data := models.GetContacts(a.DB, uint(id))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(res, resp)
}
