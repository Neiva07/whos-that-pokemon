package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

//Register save in the DB a game between 2 players
var Register = func(w http.ResponseWriter, r *http.Request) {

	newGameLog := &models.GameLog{}

	err := json.NewDecoder(r.Body).Decode(newGameLog)

	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong. Please try again later."))
	}

	response := newGameLog.Create()

	u.Response(w, response)
}

//RetrieveAllGameLogsFromUser search games from a specific user
var RetrieveAllGameLogsFromUser = func(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["id"]
	log.Println(userID)
	userGameLogs := &[]models.GameLog{}

	err := models.DB.GetDB().
		Joins("JOIN users ON users.id = game_logs.winner_id OR users.id = game_logs.loser_id").
		Where("game_logs.winner_id or game_logs.loser_id = ?", userID).
		Preload("Winner").Preload("Loser").
		Order("created_at DESC").Find(userGameLogs).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		u.Response(w, u.Message(false, "Database connection error. Try again later."))
	}

	response := u.Message(true, "All game logs from the user")

	response["gameLogs"] = userGameLogs

	u.Response(w, response)
}

//RetrieveAllGameLogsFromTwoFriends return all game of two specific friends
var RetrieveAllGameLogsFromTwoFriends = func(w http.ResponseWriter, r *http.Request) {

	// params := mux.Vars(r)
	// userID := params["id"]
	// friendID := params["friend_id"]

	// friendsGameLogs := &[]models.GameLog{}

}
