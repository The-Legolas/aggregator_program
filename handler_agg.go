package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/The-Legolas/aggregator_program/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error in converting string to time: %w", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)

	log.Printf("Collecting feeds every: %v\n", timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s, feed)
}

func scrapeFeed(s *state, feed database.Feed) {
	_, err := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, feed_data := range feedData.Channel.Item {
		PubDate, err := time.Parse(time.RFC1123Z, feed_data.PubDate)
		if err != nil {
			log.Printf("Error getting the feed Published At data: %v", err)
		}

		post_args := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       feed_data.Title,
			Url:         feed_data.Link,
			Description: sql.NullString{String: feed_data.Description, Valid: feed_data.Description != ""},
			PublishedAt: PubDate,
			FeedID:      feed.ID,
		}

		post, err := s.db.CreatePost(context.Background(), post_args)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					continue
				}
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}

		log.Printf("Saved post: %s", post.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
