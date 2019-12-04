package controllers

import (
	"net/http"
	"strconv"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

//CreateFriendship register a request to connect with the user who requested
var CreateFriendship = func(w http.ResponseWriter, r *http.Request) {

	newFriendship := &models.Friendship{}

	params := mux.Vars(r)

	us, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		u.Response(w, u.Message(false, "invalid id"))
	}
	fr, err := strconv.ParseUint(params["friend_id"], 10, 64)
	if err != nil {
		u.Response(w, u.Message(false, "invalid friend id"))
	}

	userID, friendID := uint(us), uint(fr)

	newFriendship.UserID = userID
	newFriendship.FriendID = friendID
	// json.NewDecoder(r.Body).Decode(newFriendship)

	response := newFriendship.Create()

	u.Response(w, response)

}

//AcceptRequest create a friendship when the supposedFriend accept the request
var AcceptRequest = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	us, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		u.Response(w, u.Message(false, "invalid id"))
	}
	fr, err := strconv.ParseUint(params["friend_id"], 10, 64)
	if err != nil {
		u.Response(w, u.Message(false, "invalid friend id"))
	}

	userID, friendID := uint(us), uint(fr)

	user, friend := &models.User{}, &models.User{}
	friendship := &models.Friendship{}

	err = friendship.Find(userID, friendID)
	if err != nil {
		u.Response(w, u.Message(false, "Fail to find friendship request."))
		return
	}
	if err == gorm.ErrRecordNotFound {
		u.Response(w, u.Message(false, "Request not found"))
		return
	}
	if friendship.FriendshipStatus != models.Requested {
		u.Response(w, u.Message(false, "There is no friendship request"))
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
	friendship.FriendshipStatus = models.Accepted
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
