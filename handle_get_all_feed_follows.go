package main

import (
	"GeminiZA/blogator/internal/database"
	"context"
	"net/http"
)

func (cfg *apiConfig) handleGetAllFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedsFollows, err := cfg.DB.GetAllFeedFollows(context.Background(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot get feed follows")
		return
	}
	respondWithJSON(w, http.StatusOK, feedsFollows)
}
