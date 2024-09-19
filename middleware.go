package main

import (
	"GeminiZA/blogator/internal/database"
	"context"
	"net/http"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.Contains(auth, "ApiKey ") {
			respondWithError(w, http.StatusUnauthorized, "no api key provided")
			return
		}
		apiKey := auth[7:]
		user, err := cfg.DB.GetUser(context.Background(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid api key")
			return
		}
		handler(w, r, user)
	}
}
