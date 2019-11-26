package controllers

import (
	"encoding/json"
	"net/http"
	"whos-that-pokemon/models"

	u "whos-that-pokemon/utils"
)

//SignUp create a user in the database
var SignUp = func(w http.ResponseWriter, r *http.Request) {

	newUser := &models.User{}

	err := json.NewDecoder(r.Body).Decode(newUser)

	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong. Please, try again."))
	}

	response := newUser.Create()

	u.Response(w, response)
}
