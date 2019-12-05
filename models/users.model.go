package models

import (
	"strings"

	"github.com/jinzhu/gorm"

	u "whos-that-pokemon/utils"
)

//User is a gorm model to use User
type User struct {
	gorm.Model

	GivenName  string
	FamilyName string
	ImageURL   string `json:"photo"`
	Name       string
	Email      string  `gorm:"unique;not null"`
	Token      string  `gorm:"column:token" json:"idToken"`
	GoogleID   uint    `gorm:"column:google_id"`
	GameLogs   []Game  `json:"-"`
	Friends    []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id" json:"-"`
}

//validate function check if there's not invalid in the signup process
func (user *User) validate() (map[string]interface{}, bool) {

	if !strings.Contains((user.Email), "@") {

		return u.Message(false, "Invalid email address. Please, try a real one."), false
	}

	existUser := &User{}

	err := DB.GetDB().Table("users").Where("email = ?", user.Email).First(existUser).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please try again later."), false
	}
	if existUser.Email != "" {
		return u.Message(false, "Email addres already taken."), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Create method to create user and save it in the database
func (user *User) Create() map[string]interface{} {

	if response, ok := user.validate(); ok == false {
		return response
	}

	DB.GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create an user. Connection error.")
	}

	//validate the account here

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response
}

//Find make a query to find a user in the database
func (user *User) Find(userID uint) error {

	err := DB.GetDB().Table("users").Where("id = ?", userID).First(&user).Error

	return err
}

//AssociateFriend take a friendship and associate into the type User
func (user *User) AssociateFriend(friend *User) error {

	err := DB.GetDB().Model(&user).Association("Friends").Append(friend).Error
	return err
}
