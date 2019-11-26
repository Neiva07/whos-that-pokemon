package models

import (
	"github.com/jinzhu/gorm"

	u "whos-that-pokemon/utils"
)

//User is a gorm model to use User
type User struct {
	gorm.Model

	userID     uint
	GivenName  string
	FamilyName string
	Photo      string
	Name       string
	Email      string
	Token      string
	id         uint
}

//Create method to create user and save it in the database
func (user *User) Create() map[string]interface{} {

	DB.GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create an user. Connection error.")
	}

	//validate the account here

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response
}
