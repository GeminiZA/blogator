package main

import (
	"GeminiZA/blogator/internal/database"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := r.PathValue("feedFollowID")
	fmt.Printf("Got feed follow ID: %s\n", feedFollowIDStr)
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid feed follow id")
		return
	}
	err = cfg.DB.UnfollowFeed(context.Background(), feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot unfollow feed")
		return
	}
	w.WriteHeader(http.StatusOK)
}
