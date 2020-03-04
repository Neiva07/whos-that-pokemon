package models

import (
	"encoding/json"
	"strings"
	"time"

	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"
	googleOauth "google.golang.org/api/oauth2/v1"
)

//User is a gorm model to use User
type User struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `sql:"index" json:"deletedAt"`
	GivenName  string     `json:"givenName"`
	FamilyName string     `json:"familyName"`
	ImageURL   string     `json:"photo"`
	Name       string     `json:"name"`
	Email      string     `gorm:"unique;not null" json:"email"`
	token      string     `gorm:"column:token"`
	GoogleID   string     `gorm:"column:google_id" json:"googleId"`
	GameLogs   []Game     `json:"-" json:"gameLogs"`
	Friends    []*User    `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id" json:"-"`
}

//UnmarshalBinary Overrides interface to be able to convert []byte into User model
func (user *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, user)
}

//MarshalBinary overrides interface to be able to convert User model into []byte
func (user *User) MarshalBinary() ([]byte, error) {

	return json.Marshal(user)
}

//ConvertionFromGoogleUser takes a userInfo fields from Google API and peer to correspondent database field
func (user *User) ConvertionFromGoogleUser(googleUser *googleOauth.Userinfoplus, accessToken string) {
	user.FamilyName = googleUser.FamilyName
	user.ImageURL = googleUser.Picture
	user.GivenName = googleUser.GivenName
	user.Name = googleUser.Name
	user.token = accessToken
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

//Update take a user model and changes filds in the database
func (user *User) Update() error {

	err := DB.GetDB().Model(user).Updates(user).Error

	return err
}

//Find make a query to find a user in the database
func (user *User) Find(userID uint) error {

	err := DB.GetDB().Table("users").Where("id = ?", userID).First(&user).Error

	return err
}

//FindByEmail search for a user by email address
func (user *User) FindByEmail(email string) error {
	err := DB.GetDB().Table("users").Where("email = ?", email).First(&user).Error
	return err
}

//AssociateFriend take a friendship and associate into the type User
func (user *User) AssociateFriend(friend *User) error {

	err := DB.GetDB().Model(user).Association("Friends").Append(friend).Error
	return err
}
