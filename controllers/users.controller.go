package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"whos-that-pokemon/app"
	"whos-that-pokemon/models"

	"github.com/jinzhu/gorm"

	u "whos-that-pokemon/utils"
)

//SignIn create a user in the database
var SignIn = func(w http.ResponseWriter, r *http.Request) {

	var authCode = make(map[string]string)
	var response = make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&authCode)

	if err != nil {
		response = u.Message(false, "AuthCode malformed")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Response(w, response)
		return
	}

	token, err := app.TokenHandler(authCode["authCode"], r.Context())

	if err != nil {
		response = u.Message(false, "Something went wrong creating a token")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Response(w, response)
		return
	}

	googleUser, err := app.GetGoogleUserInfo(token, r.Context())

	if err != nil {
		response = u.Message(false, "Error connection with Google's API")
		w.Header().Add("Content-Type", "application/json")
		u.Response(w, response)
		return
	}

	user := &models.User{}

	err = user.FindByEmail(googleUser.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		response = u.Message(false, "Something went wrong searching for a user in the database")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		u.Response(w, response)
		return
	}

	user.ConvertionFromGoogleUser(googleUser, token.AccessToken)

	if err == gorm.ErrRecordNotFound {
		user.GoogleID = googleUser.Id
		response = user.Create()
	} else {
		if err = user.Update(); err != nil {
			response = u.Message(false, "Error while updating user profile")
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}
	}

	err = models.Redis.Set(authCode["authCode"], user, 0).Err()

	log.Println(models.Redis.Get(authCode["authCode"]))

	if err != nil {
		response = u.Message(false, "Error in the user session system")
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		u.Response(w, response)
		return
	}

	response = u.Message(true, "User logged in!")
	response["user"] = user
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	u.Response(w, response)
	return
}

// var SignOut =
