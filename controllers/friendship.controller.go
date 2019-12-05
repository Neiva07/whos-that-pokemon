package controllers

import (
	"net/http"
	"strconv"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//CreateFriendship register a request to connect with the user who requested
var CreateFriendship = func(w http.ResponseWriter, r *http.Request) {

	newFriendship := &models.Friendship{}
	userID, friendID, err := parseUserAndFriendIds(r)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid user id or friend id."))
		return
	}

	newFriendship.UserID = userID
	newFriendship.FriendID = friendID

	response := newFriendship.Create()

	u.Response(w, response)

}

//AcceptRequest create a friendship when the Friend accept the request
var AcceptRequest = func(w http.ResponseWriter, r *http.Request) {

	user, friend := &models.User{}, &models.User{}
	friendship := &models.Friendship{}

	userID, friendID, err := parseUserAndFriendIds(r)

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
	err = friendship.Update()
	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong saving the change into the database."))
		return
	}

	err = user.AssociateFriend(friend)
	if err != nil {
		u.Response(w, u.Message(false, "Associassion error. Something went wrong creating the association."))
		return
	}

	u.Response(w, u.Message(true, "Friendship created!"))
	return
}

// DeleteFriendship delete a request or a friendship between 2 users
var DeleteFriendship = func(w http.ResponseWriter, r *http.Request) {
	friendship := &models.Friendship{}
	userID, friendID, err := parseUserAndFriendIds(r)
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

	u.Response(w, u.Message(false, "Friendship deleted"))
	return

}

func parseUserAndFriendIds(r *http.Request) (uint, uint, error) {

	params := mux.Vars(r)

	us, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	fr, err := strconv.ParseUint(params["friend_id"], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	userID, friendID := uint(us), uint(fr)

	return userID, friendID, nil
}
