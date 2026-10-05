package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func unescapeRSSItem(rssItem RSSItem) RSSItem {
	return RSSItem{
		Title:       html.UnescapeString(rssItem.Title),
		Link:        html.UnescapeString(rssItem.Link),
		Description: html.UnescapeString(rssItem.Description),
		PubDate:     html.UnescapeString(rssItem.PubDate),
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal rss feed: %v", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Link = html.UnescapeString(rssFeed.Channel.Link)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range len(rssFeed.Channel.Item) {
		rssFeed.Channel.Item[i] = unescapeRSSItem(rssFeed.Channel.Item[i])
	}
	return &rssFeed, nil
}
