package utils

import (
	"encoding/json"
	"net/http"
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
