package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"github.com/gorilla/mux"
)

//FriendshipRequest register a request to connect with the user who requested
var FriendshipRequest = func(w http.ResponseWriter, r *http.Request) {

	newFriendshipRequest := &models.FriendshipRequest{}

	json.NewDecoder(r.Body).Decode(newFriendshipRequest)

	response := newFriendshipRequest.Create()

	u.Response(w, response)

	// err = models.DB.GetDB().Table("users").Where("email = ?", email).Find(&friend).Error
	// if err == gorm.ErrRecordNotFound {
	// 	u.Response(w, u.Message(false, "Friend not found"))
	// } else if err != nil {
	// 	u.Response(w, u.Message(false, "Fail to find Friend. Fail connect with database."))
	// }
	// err = models.DB.GetDB().Table("users").Where("id = ?", UserID).Find(&user).Error
	// if err == gorm.ErrRecordNotFound {
	// 	u.Response(w, u.Message(false, "User not found"))
	// } else if err != nil {
	// 	u.Response(w, u.Message(false, "Fail to find User. Fail connect with database."))
	// }

}

//AcceptRequest create a friendship when the supposedFriend accept the request
var AcceptRequest = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	us, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		u.Response(w, u.Message(false, "invalid id"))
	}
	fr, err := strconv.ParseUint(params["friend_id"], 10, 64)

	userID, friendID := uint(us), uint(fr)

	user, friend := &models.User{}, &models.User{}
	friendshipRequest := &models.FriendshipRequest{}

	err = friendshipRequest.Find(userID, friendID)

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

	err = models.DB.GetDB().Model(&user).Association("friendship").Append(friend).Error
	if err != nil {
		u.Response(w, u.Message(false, "Associassion error. Something went wrong creating the association."))
	}

	err = friendshipRequest.Delete()

	if err != nil {
		u.Response(w, u.Message(false, "Fail to delete friendship request after accepted."))
	}

	u.Response(w, u.Message(true, "Friendship created!"))

}
