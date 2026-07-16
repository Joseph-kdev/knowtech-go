package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/Joseph-kdev/knowtech-go/response"
)

type FeedGroup struct {
	FeedID   string    `json:"feed_id"`
	FeedName string    `json:"feed_name"`
	FeedUrl  string    `json:"feed_url"`
	Posts    []RSSItem `json:"posts"`
}

type FeedPosts struct {
	FeedID   string         `json:"feed_id"`
	FeedName string         `json:"feed_name"`
	FeedUrl  string         `json:"feed_url"`
	Posts    []PostResponse `json:"posts"`
}

type PostResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt time.Time  `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func aggregatePosts(rows []db.GetPostsByFeedRow) []FeedPosts {
	feedOrder := make([]string, 0)
	feedMap := make(map[string]*FeedPosts)

	for _, row := range rows {
		feedID := row.FeedID.String()
		if _, exists := feedMap[feedID]; !exists {
			feedMap[feedID] = &FeedPosts{
				FeedID:   feedID,
				FeedName: row.FeedName,
				FeedUrl:  row.FeedUrl,
				Posts:    []PostResponse{},
			}
			feedOrder = append(feedOrder, feedID)
		}

		var desc *string
		if row.PostDescription.Valid {
			desc = &row.PostDescription.String
		}

		feedMap[feedID].Posts = append(feedMap[feedID].Posts, PostResponse{
			ID:          row.PostID.String(),
			Title:       row.PostTitle,
			Url:         row.PostUrl,
			Description: desc,
			PublishedAt: row.PostPublishedAt,
			CreatedAt:   row.PostCreatedAt,
			UpdatedAt:   row.PostUpdatedAt,
		})
	}

	result := make([]FeedPosts, 0, len(feedOrder))
	for _, id := range feedOrder {
		result = append(result, *feedMap[id])
	}
	return result
}

func aggregateFollowedPosts(rows []db.GetPostsFromFollowedFeedsRow) []FeedPosts {
	feedOrder := make([]string, 0)
	feedMap := make(map[string]*FeedPosts)

	for _, row := range rows {
		feedID := row.FeedID.String()
		if _, exists := feedMap[feedID]; !exists {
			feedMap[feedID] = &FeedPosts{
				FeedID:   feedID,
				FeedName: row.FeedName,
				FeedUrl:  row.FeedUrl,
				Posts:    []PostResponse{},
			}
			feedOrder = append(feedOrder, feedID)
		}

		var desc *string
		if row.PostDescription.Valid {
			desc = &row.PostDescription.String
		}

		feedMap[feedID].Posts = append(feedMap[feedID].Posts, PostResponse{
			ID:          row.PostID.String(),
			Title:       row.PostTitle,
			Url:         row.PostUrl,
			Description: desc,
			PublishedAt: row.PostPublishedAt,
			CreatedAt:   row.PostCreatedAt,
			UpdatedAt:   row.PostUpdatedAt,
		})
	}

	result := make([]FeedPosts, 0, len(feedOrder))
	for _, id := range feedOrder {
		result = append(result, *feedMap[id])
	}
	return result
}

func (apiCfg *Apiconfig) GetGroupedPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := apiCfg.DB.GetPostsByFeed(r.Context())
	if err != nil {
		response.RespondWithError(w, 500, "error fetching posts")
		return
	}

	response.RespondWithJSON(w, 201, aggregatePosts(rows))
}

func (apiCfg *Apiconfig) GetPostsFromFollowed(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("user_id")
	if param == "" {
		response.RespondWithError(w, 400, "user_id query parameter is required")
		return
	}

	followedPosts, err := apiCfg.DB.GetPostsFromFollowedFeeds(r.Context(), param)
	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("error fetching posts from followed feeds: %v", err))
		return
	}

	log.Printf("fetched %d rows from followed feeds", len(followedPosts))
	response.RespondWithJSON(w, 201, aggregateFollowedPosts(followedPosts))
}
