package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gundamdouble00/blog-aggregator-2/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	ctxBackground := context.Background()
	user, err := s.db.GetUser(ctxBackground, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieving user: %w", err)
	}

	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(ctxBackground, feedURL)
	if err != nil {
		return fmt.Errorf("couldn't retrieving feed: %w", err)
	}

	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed follow created:")
	fmt.Printf("\t* Feed's name: %s\n", newFeedFollow.FeedName)
	fmt.Printf("\t* Current user: %s\n", s.cfg.CurrentUserName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: following")
	}

	ctxBackground := context.Background()
	user, err := s.db.GetUser(ctxBackground, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieving user: %w", err)
	}

	followRecorlds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't retrieving feed follow: %w", err)
	}

	if len(followRecorlds) != 0 {
		fmt.Printf("Current user: %s\n", s.cfg.CurrentUserName)
		fmt.Println("Feed's name:")
		for _, feed := range followRecorlds {
			fmt.Printf("\t* %s\n", feed.FeedName)
		}
	} else {
		fmt.Println("No feed follows found for this user.")
	}

	return nil
}
