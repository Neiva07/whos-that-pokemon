package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"github.com/gorilla/mux"
)

//StartGameWithFriend handles the request to start a game beteween 2 players
var StartGameWithFriend = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		u.Response(w, u.Message(false, "Wrong user id. Please try again with a right one."))
	}
	friendID, err := strconv.ParseUint(params["friend_id"], 10, 64)

	if err != nil {
		u.Response(w, u.Message(false, "Wrong friend id. Please try again with a right one."))
	}

	newGame := &models.Game{}
	newGenerationsRecords := &[]models.Generation{}

	newGame.FriendID = uint(userID)
	newGame.UserID = uint(friendID)

	err = json.NewDecoder(r.Body).Decode(newGame)

	models.BulkCreateRecords(newGenerationsRecords)

	if err != nil {
		u.Response(w, u.Message(false, "Some data was settle wrong. Please try again later."))
	}

	response := newGame.Create()

	err = newGame.AddGenerations(newGenerationsRecords)

	if err != nil {
		u.Response(w, u.Message(false, "Fail to create relationship between game and their generations records."))
	}

	u.Response(w, response)

}
