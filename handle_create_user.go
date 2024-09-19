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

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var reqJson struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqJson)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "no 'name' provided")
		return
	}
	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      reqJson.Name,
	})
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
