package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gundamdouble00/blog-aggregator-2/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	numArgs := len(cmd.Args)
	if numArgs != 2 {
		return fmt.Errorf("usase: %v <name> <url>", cmd.Name)
	}

	ctxBackground := context.Background()
	feed, err := s.db.CreateFeed(ctxBackground, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(ctxBackground, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't saving feed follow: %w", err)
	}

	fmt.Println("Feed created successfully")
	printFieldsOfFeed(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usase: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		fmt.Printf("* Feed's name: %s\n", feed.FeedName)
		fmt.Printf("* Feed's URL: %s\n", feed.Url)
		fmt.Printf("* Created by: %s\n", feed.UserName)
		fmt.Println("#\t#\t#")
	}
	return nil
}

func printFieldsOfFeed(feed database.Feed) {
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Created At: %v\n", feed.CreatedAt)
	fmt.Printf("Update At: %v\n", feed.UpdatedAt)
	fmt.Printf("Feed's name: %v\n", feed.Name)
	fmt.Printf("Feed's URL: %v\n", feed.Url)
	fmt.Printf("Feed owner ID: %v\n", feed.UserID)
}
