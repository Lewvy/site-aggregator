package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fetchFeed: Error creating request: %q", err)
	}
	req.Header.Set("User-Agent", "gator/0.01")
	req.Header.Set("Content-Type", "application/xml")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetchFeed: Error making request: %q", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("fetchFeed: Error Reading resp: %q", err)
	}
	feed := &RSSFeed{}
	if err = xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("fetchFeed: Error unmarshalling data: %q", err)
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return feed, nil
}
