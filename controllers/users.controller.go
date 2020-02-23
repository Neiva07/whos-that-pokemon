package controllers

import (
	"encoding/json"
	"net/http"
	"whos-that-pokemon/models"

	u "whos-that-pokemon/utils"
)

//SignIn create a user in the database
var SingIn = func(w http.ResponseWriter, r *http.Request) {

	newUser := &models.User{}

	err := json.NewDecoder(r.Body).Decode(newUser)

	if err != nil {
		u.Response(w, u.Message(false, "Something went wrong. Please, try again."))
	}

	newUser = (r.Context().Value(newUser)).(*models.User)

	response := newUser.Create()

	u.Response(w, response)
}
