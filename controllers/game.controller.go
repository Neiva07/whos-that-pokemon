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

	// newGenerationsRecords, err := models.BulkCreateRecords(&newGame.Generations)

	if err != nil {
		u.Response(w, u.Message(false, "Some data was settle wrong. Please try again later."))
		return
	}

	response := newGame.Create()
	err = newGame.AddGenerations(&newGame.Generations)

	u.Response(w, response)
	return
}

//GetAllUserGames retrieve All games from a user
var GetAllUserGames = func(w http.ResponseWriter, r *http.Request) {

	userID, err := u.ParseID(r)

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

//GetSpecificGame return a specific game details to a user
var GetSpecificGame = func(w http.ResponseWriter, r *http.Request) {

	gameID, err := u.ParseID(r)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid Game ID"))
		return
	}

	game := &models.Game{}

	err = game.Find(gameID)

	if err == gorm.ErrRecordNotFound {

		u.Response(w, u.Message(false, "Game doesn't exist"))
		return
	}

	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong on searching"))
		return
	}

	response := u.Message(true, "Game found.")

	response["game"] = game

	u.Response(w, response)
	return
}

//UpdateGame change the status and the score of a game while the game is open
var UpdateGame = func(w http.ResponseWriter, r *http.Request) {

	gameID, err := u.ParseID(r)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid Game ID"))
		return
	}

	updatedGame, game := &models.Game{}, &models.Game{}

	err = json.NewDecoder(r.Body).Decode(updatedGame)

	if err != nil {
		u.Response(w, u.Message(false, "Impossible to read the data. Bad formatting"))
		return
	}

	err = game.Find(gameID)

	if err == gorm.ErrRecordNotFound {
		u.Response(w, u.Message(false, "Game not found to update"))
		return
	}
	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong connecting to the database"))
		return
	}

	err = game.Update(updatedGame)

	if err != nil {
		u.Response(w, u.Message(false, "Could not update the game. Error connecting to the database"))
		return
	}

	response := u.Message(true, "Game updated successfully!")

	response["game"] = game

	u.Response(w, response)
	return
}

// GetAllGamesFromFriends return all games for a specific friendship
var GetAllGamesFromFriends = func(w http.ResponseWriter, r *http.Request) {

	userID, friendID, err := u.ParseUserAndFriendIDs(r)
	games := &[]models.Game{}

	if err != nil {
		u.Response(w, u.Message(false, "Invalid user id or friend id format."))
		return
	}

	err = models.DB.GetDB().Table("games").
		Where("user_id = ? AND friend_id = ? OR friend_id = ? AND user_id = ?", userID, friendID, userID, friendID).
		Find(games).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		u.Response(w, u.Message(false, "Something went wrong connecting to the database"))
		return
	}

	err = models.DB.GetDB().Preload("Generations").Find(games).Error

	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong searching for generations"))
		return
	}

	response := u.Message(true, "Games found successfully!")
	response["games"] = games

	u.Response(w, response)
	return
}
