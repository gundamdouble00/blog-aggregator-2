package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func unescapeRSSItem(rssItem RSSItem) RSSItem {
	return RSSItem{
		Title:       html.UnescapeString(rssItem.Title),
		Link:        rssItem.Link,
		Description: html.UnescapeString(rssItem.Description),
		PubDate:     rssItem.PubDate,
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal rss feed: %v", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range len(rssFeed.Channel.Item) {
		rssFeed.Channel.Item[i] = unescapeRSSItem(rssFeed.Channel.Item[i])
	}
	return &rssFeed, nil
}
