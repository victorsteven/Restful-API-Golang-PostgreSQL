package models

import (
	"fmt"
	u "rest_api_gorm/utils"

	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserId uint   `json:"user_id"` //The user that this contact belongs to
}

func (contact *Contact) Validate(db *gorm.DB) (map[string]interface{}, bool) {
	if contact.Name == "" {
		return u.Message(false, "Contact name should be on the payload"), false
	}
	if contact.Phone == "" {
		return u.Message(false, "Phone number is required"), false
	}
	if contact.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}
	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (contact *Contact) Create(db *gorm.DB) map[string]interface{} {
	if resp, ok := contact.Validate(db); !ok {
		return resp
	}
	db.Create(contact)
	resp := u.Message(true, "success")
	resp["contact"] = contact
	return resp
}

// func GetContact(id uint) *Contact {
// 	contact := &Contact{}
// 	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
// 	if err != nil {
// 		return nil
// 	}
// 	return contact
// }

func GetContacts(db *gorm.DB, user uint) []*Contact {
	contacts := make([]*Contact, 0)
	err := db.Table("contacts").Where("user_id = ?", user).Find(&contacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contacts
}
