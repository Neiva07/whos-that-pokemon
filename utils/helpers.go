package utils

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Message helper function to return in a standard map form
var Message = func(status bool, s string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": s}
}

// Response add to the header and the writer the data in json format
var Response = func(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}

//ParseUserID takes the user from the param and parse to uint
var ParseUserID = func(r *http.Request) (uint, error) {

	params := mux.Vars(r)

	us, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return 0, err
	}
	userID := uint(us)

	return userID, err
}

//ParseUserAndFriendIDs takes the user and friend from the param and parse to uint
func ParseUserAndFriendIDs(r *http.Request) (uint, uint, error) {

	params := mux.Vars(r)

	us, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	fr, err := strconv.ParseUint(params["friend_id"], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	userID, friendID := uint(us), uint(fr)

	return userID, friendID, nil
}
