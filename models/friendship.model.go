package models

import (
	"log"
	"time"
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"
)

//FriendshipStatus limited the Status inside Friendship.Status
type FriendshipStatus uint

const (
	//Requested friendship
	Requested FriendshipStatus = iota + 1
	//Accepted friendship request
	Accepted
	//Deleted friendship
	Deleted
)

//Friendship is a Model to store a request to create a friendship
type Friendship struct {
	CreatedAt        time.Time        `sql:"default:now()"`
	UpdatedAt        time.Time        `sql:"default:now()"`
	DeletedAt        time.Time        `sql:"default:NULL"`
	UserID           uint             `gorm:"primary_key;auto_increment:false"`
	FriendID         uint             `gorm:"primary_key;auto_increment:false"`
	User             User             `gorm:"foreignkey:id;association_foreignkey:UserID" json:"-"`
	Friend           User             `gorm:"foreignkey:id;association_foreignkey:FriendID" json:"-"`
	FriendshipStatus FriendshipStatus `gorm:"default:1"`
}

func (friendship *Friendship) validate(userID uint, friendID uint) (map[string]interface{}, bool) {

	err := DB.GetDB().Table("users").Where("id = ?", userID).Error

	if err != nil {
		return u.Message(false, "User not found in the database."), false
	}
	err = DB.GetDB().Table("users").Where("id = ?", friendID).Error
	if err != nil {
		return u.Message(false, "User not found in the database."), false
	}

	return u.Message(true, "Request valid"), true

}

//Create a new friendship request in the database
func (friendship *Friendship) Create(userID uint, friendID uint) map[string]interface{} {

	if msg, ok := friendship.validate(userID, friendID); !ok {
		return msg
	}

	err := DB.GetDB().Table("friendships").Unscoped().
		Where("user_id = ? AND friend_id = ? OR user_id = ? AND friend_id = ?", userID, friendID, friendID, userID).
		Attrs(Friendship{UserID: userID, FriendID: friendID}).
		FirstOrCreate(friendship).Error

	if friendship.FriendshipStatus == Deleted {

		err = friendship.Update(&Friendship{FriendshipStatus: Requested})
		if err != nil {
			return u.Message(false, "Could not update friendship after deleted.")
		}
	}

	if friendship.FriendshipStatus == Accepted {

		return u.Message(false, "Friendship already exist")
	}

	if err != nil {
		return u.Message(false, "Could not create friendship request.")
	}

	response := u.Message(true, "Friendship Request created successfully")

	response["friendship"] = friendship

	return response
}

//Find a specific friendship request between 2 users
func (friendship *Friendship) Find(userID uint, friendID uint) error {

	err := DB.GetDB().Table("friendships").Where("user_id = ? AND friend_id = ?", userID, friendID).Find(friendship).Error
	log.Print(friendship.FriendshipStatus)
	return err
}

//Delete delete a friendship request
func (friendship *Friendship) Delete() error {
	err := DB.GetDB().Model(friendship).Update("friendship_status", Deleted).Error
	if err != nil {
		return err
	}
	err = DB.GetDB().Table("friendships").Delete(friendship).Error
	return err
}

//Update saves the current updated friendship instance
func (friendship *Friendship) Update(FriendshipFields *Friendship) error {
	err := DB.GetDB().Unscoped().
		Model(friendship).Updates(FriendshipFields).Error
	log.Println(FriendshipFields.FriendshipStatus)
	return err
}

//BeforeUpdate a deleted friendship, the delete_at is set to nil
func (friendship *Friendship) BeforeUpdate(scope *gorm.Scope) {
	if friendship.FriendshipStatus != Deleted {
		scope.SetColumn("deleted_at", nil)
	}
}
