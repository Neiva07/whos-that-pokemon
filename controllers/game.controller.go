package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"
)

//StartGameWithFriend handles the request to start a game beteween 2 players
var StartGameWithFriend = func(w http.ResponseWriter, r *http.Request) {

	userID, friendID, err := u.ParseUserAndFriendIDs(r)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid user id or friend id."))
		return
	}

	newGame := &models.Game{}

	newGame.FriendID = uint(userID)
	newGame.UserID = uint(friendID)

	err = json.NewDecoder(r.Body).Decode(newGame)

	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong decoding the json."))
		return
	}

	newGenerationsRecords, err := models.BulkCreateRecords(&newGame.Generations)

	if err != nil {
		u.Response(w, u.Message(false, "Some data was settle wrong. Please try again later."))
		return
	}

	response := newGame.Create()
	err = newGame.AddGenerations(newGenerationsRecords)

	u.Response(w, response)
	return
}

//GetAllUserGames retrieve All games from a user
var GetAllUserGames = func(w http.ResponseWriter, r *http.Request) {

	userID, err := u.ParseUserID(r)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid user id."))
		return
	}

	var results []struct {
		GivenName   string
		FamilyName  string
		Name        string
		Email       string
		UserScore   uint
		FriendScore uint
		FriendID    uint
		UserID      uint
		CreatedAt   time.Time
		Timer       uint
		Status      models.GameStatus
	}

	selection := `
			users.given_name,
			users.family_name,
			users.name,
			users.email,
			games.user_score,
			games.friend_score,
			games.friend_id,
			games.user_id,
			games.created_at,
			games.status,
			games.timer
	`

	err = models.DB.GetDB().Table("games").Select(selection).
		Joins("JOIN users ON users.id = games.user_id AND games.friend_id = ?", userID).
		Joins("UNION ?", models.DB.GetDB().Table("games").Select(selection).QueryExpr()).
		Joins("JOIN users ON users.id = games.friend_id AND games.user_id = ?", userID).
		Scan(&results).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		u.Response(w, u.Message(false, "Something went wrong connecting with the database"))
		return
	}

	response := u.Message(true, "Games found successfully!")
	response["games"] = results

	u.Response(w, response)
	return

}
