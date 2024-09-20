package main

import (
	"GeminiZA/blogator/internal/database"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var reqBody struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "feed id not provided in json body")
		return
	}
	feedFollow, err := cfg.DB.FollowFeed(context.Background(), database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    reqBody.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot follow feed")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}
