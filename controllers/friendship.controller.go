package controllers

import (
	"net/http"
	"time"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"
)

//CreateFriendship register a request to connect with the user who requested
var CreateFriendship = func(w http.ResponseWriter, r *http.Request) {

	newFriendship := &models.Friendship{}
	userID, friendID, err := u.ParseUserAndFriendIDs(r)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid user id or friend id."))
		return
	}

	response := newFriendship.Create(userID, friendID)

	u.Response(w, response)

}

// SearchAllFriends return all friends from a user
var SearchAllFriends = func(w http.ResponseWriter, r *http.Request) {

	userID, err := u.ParseID(r)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid user id."))
		return
	}

	var results []struct {
		GivenName        string
		FamilyName       string
		Name             string
		UserID           uint
		Email            string
		ImageURL         string
		FriendshipStatus uint
		UserTotalScore   uint
		FriendTotalScore uint
		CreatedAt        time.Time
	}

	selection := `users.given_name, 
				users.family_name,
				users.id as user_id, 
				users.email, 
				users.name, 
				users.image_url, 
				friendships.friendship_status, 
				friendships.user_total_score,
				friendships.friend_total_score
				friendships.created_at`

	err = models.DB.GetDB().Table("friendships").Select(selection).
		Joins("JOIN users ON users.id = friendships.user_id AND  friendships.deleted_at IS NULL AND friendships.friend_id = ?", userID).
		Joins("UNION ?", models.DB.GetDB().Table("friendships").Select(selection).QueryExpr()).
		Joins("JOIN users ON users.id = friendships.friend_id AND friendships.deleted_at IS NULL AND friendships.user_id = ? AND friendships.friendship_status = 2", userID).
		Scan(&results).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		u.Response(w, u.Message(false, "Something went wrong querying the data."))
		return
	}

	response := u.Message(true, "Friendships found successfully")

	response["friendships"] = results
	u.Response(w, response)
	return

}

//AcceptRequest create a friendship when the Friend accept the request
var AcceptRequest = func(w http.ResponseWriter, r *http.Request) {

	user, friend := &models.User{}, &models.User{}
	friendship := &models.Friendship{}

	userID, friendID, err := u.ParseUserAndFriendIDs(r)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid user id or friend id."))
		return
	}

	err = friendship.Find(userID, friendID)

	if err == gorm.ErrRecordNotFound {
		u.Response(w, u.Message(false, "Request not found"))
		return
	}
	if err != nil {
		u.Response(w, u.Message(false, "Fail to connect with to the database."))
		return
	}
	if friendship.FriendshipStatus != models.Requested {
		u.Response(w, u.Message(false, "There is no friendship request"))
		return
	}

	err = user.Find(userID)
	if err != nil {
		u.Response(w, u.Message(false, "Fail to find user in the database"))
		return
	}
	err = friend.Find(friendID)
	if err != nil {
		u.Response(w, u.Message(false, "Fail to find friend in the database"))
		return
	}
	err = friendship.Update(&models.Friendship{FriendshipStatus: models.Accepted})
	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong saving the change into the database."))
		return
	}

	u.Response(w, u.Message(true, "Friendship created!"))
	return
}

// DeleteFriendship delete a request or a friendship between 2 users
var DeleteFriendship = func(w http.ResponseWriter, r *http.Request) {
	friendship := &models.Friendship{}
	userID, friendID, err := u.ParseUserAndFriendIDs(r)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid user id or friend id."))
		return
	}

	err = friendship.Find(userID, friendID)
	if err != nil {
		u.Response(w, u.Message(false, "Friendship not found."))
		return
	}

	err = friendship.Delete()

	if err != nil {
		u.Response(w, u.Message(false, "Could not delete the friendship"))
		return
	}

	u.Response(w, u.Message(true, "Friendship deleted"))
	return

}
