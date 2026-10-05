package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gundamdouble00/blog-aggregator-2/internal/database"
)

func printFieldsOfFeed(feed database.Feed) {
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Created At: %v\n", feed.CreatedAt)
	fmt.Printf("Update At: %v\n", feed.UpdatedAt)
	fmt.Printf("Feed's name: %v\n", feed.Name)
	fmt.Printf("Feed's URL: %v\n", feed.Url)
	fmt.Printf("Feed owner ID: %v\n", feed.UserID)
}

func handlerAddFeed(s *state, cmd command) error {
	numArgs := len(cmd.Args)
	if numArgs != 2 {
		return fmt.Errorf("usase: %v <name> <url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieve user: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
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

	fmt.Println("Feed created successfully")
	printFieldsOfFeed(feed)
	return nil
}
