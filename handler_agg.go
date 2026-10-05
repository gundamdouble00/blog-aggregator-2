package main

import (
	"context"
	"fmt"
)

const RSS_URL = "https://www.wagslane.dev/index.xml"
const MDN_RSS_URL = "https://developer.mozilla.org/en-US/blog/rss.xml"
const MEDIUM = "https://medium.com"

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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("\"agg\" command doesn't have any arguments")
	}

	rssFeed, err := fetchFeed(context.Background(), RSS_URL)
	if err != nil {
		return err
	}

	printRSSFeed(*rssFeed)
	return nil
}
