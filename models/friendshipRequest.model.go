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

//Find a specific friendship request between 2 users
func (friendshipRequest *FriendshipRequest) Find(userID uint, friendID uint) error {

	err := DB.GetDB().Table("friendship_requests").Where("user_id = ? AND friend_id = ?", userID, friendID).First(friendshipRequest).Error

	return err
}

//Delete delete a friendship request
func (friendshipRequest *FriendshipRequest) Delete() error {

	err := DB.GetDB().Table("friendship_requests").Delete(friendshipRequest).Error

	return err
}
