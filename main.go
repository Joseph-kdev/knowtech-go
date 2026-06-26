package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Joseph-kdev/knowtech-go/handlers"
	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
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
	
	apiCfg := handlers.Apiconfig{DB: sqlDB}

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

	apiRouter.Post("/feed", apiCfg.AddFeed)
	apiRouter.Get("/feed", apiCfg.GetAllFeeds)

	apiRouter.Get("/posts", apiCfg.GetGroupedPosts)

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
 