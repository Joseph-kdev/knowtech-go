package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Joseph-kdev/knowtech-go/response"
)

type FeedGroup struct {
	FeedID   string    `json:"feed_id"`
	FeedName string    `json:"feed_name"`
	FeedUrl  string    `json:"feed_url"`
	Posts    []RSSItem `json:"posts"`
}

func (apiCfg *Apiconfig) GetGroupedPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := apiCfg.DB.GetPostsByFeed(r.Context())
	if err != nil {
		response.RespondWithError(w, 500, "error fetching posts")
		return
	}

	result := make([]FeedGroup, 0, len(rows))
	for _, row := range rows {
		var posts []RSSItem
		if err := json.Unmarshal(row.Posts, &posts); err != nil {
			log.Printf("error unmarshalling posts for feed %s: %v", row.FeedName, err)
			continue
		}
		result = append(result, FeedGroup{
			FeedID:   row.FeedID.String(),
			FeedName: row.FeedName,
			FeedUrl:  row.FeedUrl,
			Posts:    posts,
		})
	}
	response.RespondWithJSON(w, 201, result)
}