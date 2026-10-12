package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
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

	printRSSFeed(*rssFeed)
	return nil
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
