package models

import (
	"os"
	u "rest_api_gorm/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//Account struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

//Validate incoming user details
func (account *Account) Validate(db *gorm.DB) (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email Address is required"), false
	}
	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}
	//Email must be unique
	temp := &Account{}

	//check for errros and duplicate emails
	err := db.Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please try again"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use"), false
	}
	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create(db *gorm.DB) map[string]interface{} {
	if resp, ok := account.Validate(db); !ok {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	db.Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account")
	}
	//Create a new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(db *gorm.DB, email, password string) map[string]interface{} {
	account := &Account{}
	err := db.Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email not found")
		}
		return u.Message(false, "Connection error. Try again")
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials.")
	}
	account.Password = ""

	//create JWT token:
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

func GetUser(db *gorm.DB, u uint) *Account {
	acc := &Account{}
	db.Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //Email not found
		return nil
	}

	acc.Password = "" //dont return the password with the response

	return acc
}
