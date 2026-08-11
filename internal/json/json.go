package json

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteData(w http.ResponseWriter, r *http.Request, data any) bool {
	// log.Println("inside write data and data is : ", json.NewDecoder(r.Body))
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		Write(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return false
	}
	return true
}
