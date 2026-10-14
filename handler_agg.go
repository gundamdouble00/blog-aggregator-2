package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gundamdouble00/blog-aggregator-2/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolationErr = "23505"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v...\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			log.Print(err)
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("not following any feeds")
		}

		return fmt.Errorf("couldn't get next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("couldn't mark feed fetch: %w", err)
	}

	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	fmt.Printf("Fetch feed %v (%v)\n", rssFeed.Channel.Title, rssFeed.Channel.Link)
	for _, postInfo := range rssFeed.Channel.Item {
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     postInfo.Title,
			Url:       postInfo.Link,
			Description: sql.NullString{
				String: postInfo.Description,
				Valid:  true,
			},
			PublishedAt: parsePublishedAt(postInfo.PubDate),
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationErr {
				continue
			}

			log.Println(err)
		}
	}

	// printRSSFeed(*rssFeed)
	return nil
}

func parsePublishedAt(date string) sql.NullTime {
	date = strings.TrimSpace(date)
	layouts := []string{
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, date)
		if err == nil {
			return sql.NullTime{
				Time:  parsed,
				Valid: true,
			}
		}
	}
	return sql.NullTime{}
}

func printRSSFeed(rssFeed RSSFeed) {
	channel := rssFeed.Channel
	fmt.Printf("Title: %s\n", channel.Title)
	fmt.Printf("Link: %s\n", channel.Link)
	fmt.Printf("Description: %s\n", channel.Description)
	fmt.Println("Item:")
	for _, item := range channel.Item {
		fmt.Printf("\tTitle: %s\n", item.Title)
		fmt.Printf("\tLink: %s\n", item.Link)
		fmt.Printf("\tDescription: %s\n", item.Description)
		fmt.Printf("\tPublic Date: %s\n", item.PubDate)
		fmt.Println()
	}
}
