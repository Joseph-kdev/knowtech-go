package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/Joseph-kdev/knowtech-go/response"
	"github.com/google/uuid"
)

type Apiconfig struct {
	DB *db.Queries
}

func (apiCfg *Apiconfig) AddFeed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Url      string `json:"url"`
		Category string `json:"category"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	if params.Name == "" || params.Url == "" {
		response.RespondWithError(w, 400, "name and url are required fields")
		return
	}

	exists, err := apiCfg.DB.FeedExists(r.Context(), params.Url)
	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("error checking feed existence: %v", err))
		return
	}

	if exists {
		response.RespondWithError(w, 400, "feed already exists")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID: uuid.New().String(),
		Name: params.Name,
		Url: params.Url,
		Category: sql.NullString{
			String: params.Category,
			Valid:  true,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("error creating feed: %v", err))
		return
	}

	response.RespondWithJSON(w, 201, feed)
}

func (apiCfg *Apiconfig) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("error fetching feeds: %v", err))
		return
	}
	response.RespondWithJSON(w, 200, feeds)
}