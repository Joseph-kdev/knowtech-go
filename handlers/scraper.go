package handlers

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/google/uuid"
)

func StartScraper(queries *db.Queries, concurrency int, interval time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, interval)
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		feeds, err := queries.GetFeedstoFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(queries, feed, wg)
		}
		wg.Wait()
	}
}

func scrapeFeed(queries *db.Queries, feed db.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := queries.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
	}

	rssFeed, err := FetchRSSFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching RSS feed for %s: %v", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing publication date for %s: %v", item.Link, err)
			continue
		}
		_, err = queries.AddPostsToDatabase(context.Background(), db.AddPostsToDatabaseParams{
			ID: uuid.New(),
			FeedID: feed.ID,
			Title: item.Title,
			Url: item.Link,
			Description: description,
			PublishedAt: pubAt,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error adding post to database for %s: %v", item.Link, err)
			continue
		}
		log.Printf("feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
	}
}

