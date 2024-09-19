package main

import (
	"GeminiZA/blogator/internal/database"
	"net/http"
)

func (cfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, user)
}
