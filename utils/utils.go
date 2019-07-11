package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func Respond(res http.ResponseWriter, data map[string]interface{}) {
	res.Header().Add("Content-Type", "application/json")
	json.NewEncoder(res).Encode(data)
}
