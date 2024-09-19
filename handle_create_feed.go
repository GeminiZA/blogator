package main

import (
	"GeminiZA/blogator/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var reqBody struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed json body")
		return
	}
	feed, err := cfg.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      reqBody.Name,
		Url:       reqBody.URL,
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Printf("Error creating feed: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "cannot create feed")
		return
	}
	feedFollow, err := cfg.DB.FollowFeed(context.Background(), database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	resJson := struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:       feed,
		FeedFollow: feedFollow,
	}
	respondWithJSON(w, http.StatusOK, resJson)
}
