package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Joseph-kdev/knowtech-go/handlers"
	"github.com/Joseph-kdev/knowtech-go/internal/db"
	authmw "github.com/Joseph-kdev/knowtech-go/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"google.golang.org/genai"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	db_URL := os.Getenv("DB_CONNECTION")
	if db_URL == "" {
		log.Fatal("Database connection url not found in environment")
	}

	conn, err := sql.Open("postgres", db_URL)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB := db.New(conn)

	geminiKey := os.Getenv("GEMINI_KEY")
	if geminiKey == "" {
		log.Fatal("Gemini Api Key not found in environment")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := handlers.Apiconfig{DB: sqlDB, Gemini: client}

	go handlers.StartScraper(sqlDB, 5, 4*time.Hour)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	apiRouter := chi.NewRouter()
	apiRouter.Get("/health", handlers.HandlerReadiness)

	apiRouter.Group(func(r chi.Router) {
		r.Use(authmw.MiddlewareAuth)
		r.Post("/feeds", apiCfg.AddFeed)
		r.Get("/feeds", apiCfg.GetAllFeeds)
		r.Get("/posts", apiCfg.GetPostsFromFollowed)
		r.Post("/users", apiCfg.AddUser)
		r.Post("/follow_feeds", apiCfg.FollowFeedAsUser)
		r.Get("/followed_feeds", apiCfg.GetUserFollowedFeeds)
		r.Post("/unfollow_feeds", apiCfg.UnfollowFeedAsUser)
		r.Post("/bookmarks", apiCfg.AddBookmark)
		r.Get("/bookmarks", apiCfg.GetBookmarks)
		r.Post("/remove-bookmarks", apiCfg.DeleteBookmarkForUser)
		r.Post("/chats", apiCfg.HandleChats)
	})
	apiRouter.Post("/seed", apiCfg.SeedFeeds)

	router.Mount("/api", apiRouter)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
