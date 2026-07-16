package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/Joseph-kdev/knowtech-go/response"
	"github.com/google/uuid"
)

func (apiCfg *Apiconfig) FollowFeedAsUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID string `json:"user_id"`
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	feedFollowed, err := apiCfg.DB.FollowFeed(r.Context(), db.FollowFeedParams{
		ID:        uuid.New(),
		UserID:    params.UserID,
		FeedID:    params.FeedID,
		CreatedAt: time.Now().UTC(),
	})

	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("Error following feed: %v", err))
	}

	_, err = apiCfg.DB.IncrementFeedFollowerCount(r.Context(), params.FeedID)
	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("Error incrementing feed follower count: %v", err))
		return
	}

	response.RespondWithJSON(w, 200, feedFollowed)
}

func (apiCfg *Apiconfig) GetUserFollowedFeeds(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("user_id")
	if param == "" {
		response.RespondWithError(w, 400, "user_id query parameter is required")
		return
	}

	feeds, err := apiCfg.DB.GetFollowedFeedsByUser(r.Context(), param)
	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("Error fetching followed feeds: %v", err))
		return
	}

	response.RespondWithJSON(w, 200, response.FormatFollowedFeeds(feeds))
}

func (apiCfg *Apiconfig) UnfollowFeedAsUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID string `json:"user_id"`
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	_, err = apiCfg.DB.UnfollowFeed(r.Context(), db.UnfollowFeedParams{
		FeedID: params.FeedID,
		UserID: params.UserID,
	})

	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("Error unfollowing feed: %v", err))
		return
	}

	_, err = apiCfg.DB.DecrementFeedFollowerCount(r.Context(), params.FeedID)
	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("Error decrementing feed follower count: %v", err))
		return
	}
	response.RespondWithJSON(w, 200, "Successfully unfollowed Feed")
}
