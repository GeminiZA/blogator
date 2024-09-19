package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	if code == http.StatusOK {
		fmt.Printf("Responding with JSON (%d: %s)\n", code, http.StatusText(code))
	}
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("error marshalling json response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	fmt.Printf("Responding with Error (%d: %s); error: %s\n", code, http.StatusText(code), msg)
	respondWithJSON(w, code, struct {
		Error string `json:"error"`
	}{
		Error: msg,
	})
}
