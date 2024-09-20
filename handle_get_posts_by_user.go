package main

import (
	"GeminiZA/blogator/internal/database"
	"context"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handleGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	postsLimitStr := r.URL.Query().Get("limit")
	postsLimit := 0
	var err error
	if postsLimitStr != "" {
		postsLimit, err = strconv.Atoi(postsLimitStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid limit query parameter")
			return
		}
	}
	if postsLimit == 0 || postsLimit > 500 {
		postsLimit = 500
	}
	posts, err := cfg.DB.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(postsLimit),
	})
	ret := make([]Post, len(posts))
	for i := range ret {
		ret[i] = databasePostToPost(posts[i])
	}
	respondWithJSON(w, http.StatusOK, ret)
}
