package controllers

import (
	"log"
	"net/http"
	"strconv"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

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
	}

	err = user.Find(userID)
	if err != nil {
		u.Response(w, u.Message(false, "Fail to find user in the database"))
	}
	err = friend.Find(friendID)
	if err != nil {
		u.Response(w, u.Message(false, "Fail to find friend in the database"))
	}

	err = models.DB.GetDB().Model(&user).Association("Friends").Append(friend).Error
	if err != nil {
		log.Println(err)
		u.Response(w, u.Message(false, "Associassion error. Something went wrong creating the association."))
	}

	err = friendship.Delete()

	if err != nil {
		u.Response(w, u.Message(false, "Fail to delete friendship request after accepted."))
	}

	u.Response(w, u.Message(true, "Friendship created!"))

}
