package main

import (
	"GeminiZA/blogator/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	postgresUrl := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	cfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthz", handleReadiness)
	mux.HandleFunc("/v1/err", handleError)
	mux.HandleFunc("POST /v1/users", cfg.handleCreateUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.handleGetUser))
	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.handleCreateFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.handleGetAllFeeds)
	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.handleFollowFeed))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.handleUnfollowFeed))
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.handleGetAllFeedFollows))

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
