package controllers

import (
	"encoding/json"
	"net/http"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"
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
