package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handleGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get feeds")
		return
	}
	ret := make([]Feed, len(feeds))
	for i := range ret {
		ret[i] = databaseFeedToFeed(feeds[i])
	}
	respondWithJSON(w, http.StatusOK, feeds)
}
