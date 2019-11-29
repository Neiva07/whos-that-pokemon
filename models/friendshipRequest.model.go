package models

import (
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"
)

//FriendshipRequest is a Model to store a request to create a friendship
type FriendshipRequest struct {
	gorm.Model

	UserRequestedID  uint
	SupposedFriendID uint
	User             User `gorm:"foreignkey:ID;association_foreignkey:UserRequestedID" json:"-"`
	SupposedFriend   User `gorm:"foreignkey:ID;association_foreignkey:SupposedFriendID" json:"-"`
}

func (friendshipRequest *FriendshipRequest) validate() (map[string]interface{}, bool) {

	err := DB.GetDB().Table("users").Where("id = ?", friendshipRequest.UserRequestedID).Error

	if err != nil {
		return u.Message(false, "User not found in the database."), false
	}
	err = DB.GetDB().Table("users").Where("id = ?", friendshipRequest.SupposedFriendID).Error
	if err != nil {
		return u.Message(false, "User not found in the database."), false
	}

	return u.Message(true, "Request valid"), true

}

//Create a new friendship request in the database
func (friendshipRequest *FriendshipRequest) Create() map[string]interface{} {

	if msg, ok := friendshipRequest.validate(); !ok {
		return msg
	}

	DB.GetDB().Create(friendshipRequest)

	response := u.Message(true, "Friendship Request created successfully")

	response["friendshipRequest"] = friendshipRequest

	return response
}
