package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/Joseph-kdev/knowtech-go/response"
	"github.com/google/uuid"
)

func (apiCfg *Apiconfig) AddBookmark(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string     `json:"title"`
		Url         string     `json:"url"`
		Description string     `json:"description"`
		UserID      string     `json:"user_id"`
		FeedID      uuid.UUID  `json:"feed_id"`
		FeedName    string     `json:"feed_name"`
		PublishedAt CustomTime `json:"published_at"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	if params.UserID == "" {
		response.RespondWithError(w, 400, "important details missing")
		return
	}
	insertedBookmark, err := apiCfg.DB.AddToBookmarks(r.Context(), db.AddToBookmarksParams{
		ID:    uuid.New(),
		Title: params.Title,
		Url:   params.Url,
		Description: sql.NullString{
			String: params.Description,
			Valid:  true,
		},
		UserID:      params.UserID,
		FeedID:      params.FeedID,
		PublishedAt: params.PublishedAt.Time,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})

	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("error adding bookmark: %v", err))
		return
	}

	response.RespondWithJSON(w, 200, response.FormatBookmark(insertedBookmark, params.FeedName))
}

func (apiCfg *Apiconfig) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("user_id")
	if param == "" {
		response.RespondWithError(w, 400, "user_id query parameter is required")
		return
	}

	bookmarks, err := apiCfg.DB.GetBookmarksByUser(r.Context(), param)
	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("error fetching bookmarks: %v", err))
		return
	}
	response.RespondWithJSON(w, 200, response.FormatBookmarks(bookmarks, param))
}

func (apiCfg *Apiconfig) DeleteBookmarkForUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID string    `json:"user_id"`
		Id     uuid.UUID `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	if params.UserID == "" {
		response.RespondWithError(w, 400, "user_id is missing")
		return
	}

	_, err = apiCfg.DB.DeleteBookmark(r.Context(), db.DeleteBookmarkParams{
		UserID: params.UserID,
		ID:     params.Id,
	})

	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("Error removing bookmark: %v", err))
		return
	}

	response.RespondWithJSON(w, 200, map[string]bool{"success": true})
}

type CustomTime struct {
    time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
    s := strings.Trim(string(b), "\"")
    // Try multiple formats
    formats := []string{
        "2006-01-02T15:04:05",
        "2006-01-02T15:04:05Z",
        "2006-01-02T15:04:05Z07:00",
        "2006-01-02T15:04:05.999Z07:00",
    }
    
    for _, format := range formats {
        t, err := time.Parse(format, s)
        if err == nil {
            ct.Time = t
            return nil
        }
    }
    return fmt.Errorf("unable to parse time: %s", s)
}