package models

import (
	"log"
	"time"
	u "whos-that-pokemon/utils"
)

//FriendshipStatus limited the Status inside Friendship.Status
type FriendshipStatus int

const (
	//Requested friendship
	Requested FriendshipStatus = iota
	//Accepted friendship request
	Accepted
	//Deleted friendship
	Deleted
)

//Friendship is a Model to store a request to create a friendship
type Friendship struct {
	CreatedAt        time.Time `sql:"default:now()"`
	UpdatedAt        time.Time `sql:"default:now()"`
	DeletedAt        time.Time `sql:"default:NULL"`
	UserID           uint      `gorm:"primary_key;auto_increment:false"`
	FriendID         uint      `gorm:"primary_key;auto_increment:false"`
	User             User      `gorm:"foreignkey:id;association_foreignkey:UserID" json:"-"`
	Friend           User      `gorm:"foreignkey:id;association_foreignkey:FriendID" json:"-"`
	FriendshipStatus FriendshipStatus
}

func (friendship *Friendship) validate() (map[string]interface{}, bool) {

	err := DB.GetDB().Table("users").Where("id = ?", friendship.UserID).Error

	if err != nil {
		return u.Message(false, "User not found in the database."), false
	}
	err = DB.GetDB().Table("users").Where("id = ?", friendship.FriendID).Error
	if err != nil {
		return u.Message(false, "User not found in the database."), false
	}

	// err = DB.GetDB().Table("friendship_requests").Where("user_requested_id = ? AND supposed_friend_id = ? OR supposed_friend_id = ? AND user_requested_id = ?", friendshipRequest.UserID, friendshipRequest.FriendID, friendshipRequest.FriendID, friendshipRequest.UserID).Error

	// if err != gorm.ErrRecordNotFound {
	// 	return u.Message(false, "Friendship request ")
	// }

	return u.Message(true, "Request valid"), true

}

//Create a new friendship request in the database
func (friendship *Friendship) Create() map[string]interface{} {

	if msg, ok := friendship.validate(); !ok {
		return msg
	}

	err := DB.GetDB().Create(friendship).Error

	if err != nil {
		return u.Message(false, "Could not create friendship request.")
	}

	response := u.Message(true, "Friendship Request created successfully")

	response["friendship"] = friendship

	return response
}

//Find a specific friendship request between 2 users
func (friendship *Friendship) Find(userID uint, friendID uint) error {

	err := DB.GetDB().Table("friendships").Where("user_requested_id = ? AND supposed_friend_id = ?", userID, friendID).Scan(friendship).Error
	log.Println(err)
	return err
}

//Delete delete a friendship request
func (friendship *Friendship) Delete() error {

	err := DB.GetDB().Table("friendships").Delete(friendship).Error
	return err
}
