package response

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/google/uuid"
)

func TestFormatFeed(t *testing.T) {
	feedID := uuid.New()
	category := "tech"
	lastFetchedAt := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	followersCount := int32(12)

	feed := db.Feed{
		ID:                 feedID,
		Name:               "Go News",
		Url:                "https://example.com/rss",
		Category:           sql.NullString{String: category, Valid: true},
		CreatedAt:          time.Date(2024, time.January, 1, 1, 2, 3, 0, time.UTC),
		UpdatedAt:          time.Date(2024, time.January, 1, 1, 2, 3, 0, time.UTC),
		LastFetchedAt:      sql.NullTime{Time: lastFetchedAt, Valid: true},
		FeedFollowersCount: sql.NullInt32{Int32: followersCount, Valid: true},
	}

	formatted := FormatFeed(feed)

	if formatted.ID != feedID.String() {
		t.Fatalf("expected ID %q, got %q", feedID.String(), formatted.ID)
	}
	if formatted.URL != feed.Url {
		t.Fatalf("expected URL %q, got %q", feed.Url, formatted.URL)
	}
	if formatted.Category == nil || *formatted.Category != category {
		t.Fatalf("expected category %q, got %#v", category, formatted.Category)
	}
	if formatted.LastFetchedAt == nil || !formatted.LastFetchedAt.Equal(lastFetchedAt) {
		t.Fatalf("expected last fetched time %v, got %#v", lastFetchedAt, formatted.LastFetchedAt)
	}
	if formatted.FeedFollowersCount == nil || *formatted.FeedFollowersCount != followersCount {
		t.Fatalf("expected follower count %d, got %#v", followersCount, formatted.FeedFollowersCount)
	}
}
