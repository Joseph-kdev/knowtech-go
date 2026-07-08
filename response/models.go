package response

import (
	"database/sql"
	"time"

	"github.com/Joseph-kdev/knowtech-go/internal/db"
)

type Feed struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	URL                string     `json:"url"`
	Category           *string    `json:"category,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	LastFetchedAt      *time.Time `json:"last_fetched_at,omitempty"`
	FeedFollowersCount *int32     `json:"feed_followers_count,omitempty"`
}

func FormatFeed(feed db.Feed) Feed {
	formatted := Feed{
		ID:        feed.ID.String(),
		Name:      feed.Name,
		URL:       feed.Url,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
	}

	if feed.Category.Valid {
		category := feed.Category.String
		formatted.Category = &category
	}

	if feed.LastFetchedAt.Valid {
		lastFetchedAt := feed.LastFetchedAt.Time
		formatted.LastFetchedAt = &lastFetchedAt
	}

	if feed.FeedFollowersCount.Valid {
		followersCount := feed.FeedFollowersCount.Int32
		formatted.FeedFollowersCount = &followersCount
	}

	return formatted
}

func FormatFollowedFeed(feed db.GetFollowedFeedsByUserRow) Feed {
	formatted := Feed{}

	if feed.ID.Valid {
		formatted.ID = feed.ID.UUID.String()
	}
	if feed.Name.Valid {
		formatted.Name = feed.Name.String
	}
	if feed.Url.Valid {
		formatted.URL = feed.Url.String
	}
	if feed.Category.Valid {
		category := feed.Category.String
		formatted.Category = &category
	}
	if feed.CreatedAt.Valid {
		formatted.CreatedAt = feed.CreatedAt.Time
	}
	if feed.UpdatedAt.Valid {
		formatted.UpdatedAt = feed.UpdatedAt.Time
	}
	if feed.LastFetchedAt.Valid {
		lastFetchedAt := feed.LastFetchedAt.Time
		formatted.LastFetchedAt = &lastFetchedAt
	}
	if feed.FeedFollowersCount.Valid {
		followersCount := feed.FeedFollowersCount.Int32
		formatted.FeedFollowersCount = &followersCount
	}

	return formatted
}

func FormatFeeds(feeds []db.Feed) []Feed {
	formattedFeeds := make([]Feed, 0, len(feeds))
	for _, feed := range feeds {
		formattedFeeds = append(formattedFeeds, FormatFeed(feed))
	}
	return formattedFeeds
}

func FormatFollowedFeeds(feeds []db.GetFollowedFeedsByUserRow) []Feed {
	formattedFeeds := make([]Feed, 0, len(feeds))
	for _, feed := range feeds {
		formattedFeeds = append(formattedFeeds, FormatFollowedFeed(feed))
	}
	return formattedFeeds
}

func formatNullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	formatted := value.String
	return &formatted
}
