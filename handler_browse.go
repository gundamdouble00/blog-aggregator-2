package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gundamdouble00/blog-aggregator-2/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %v <limit> | <>", cmd.Name)
	}

	limit := 2
	if len(cmd.Args) == 1 {
		num, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("limit must be a positive whole number")
		}

		if num <= 0 {
			return fmt.Errorf("limit must be a positive whole number")
		}

		limit = num
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("\t* Title: %s\n", post.Title)
		fmt.Printf("\t* URL: %s\n", post.Url)
		fmt.Printf("\t* Description: %v\n", post.Description)
		fmt.Printf("\t* Published at: %v\n", post.PublishedAt)
	}
	return nil
}
