package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Joseph-kdev/knowtech-go/handlers"
	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_"github.com/mattn/go-sqlite3"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	conn, err := sql.Open("sqlite3", "./feeds.db?_journal=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatal(err)
	}
	sqlDB := db.New(conn)
	
	apiCfg := handlers.Apiconfig{DB: sqlDB}

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

	apiRouter.Post("/feed", apiCfg.AddFeed)
	apiRouter.Get("/feed", apiCfg.GetAllFeeds)

	router.Mount("/api", apiRouter)
	
	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
