package main

import (
	"GeminiZA/blogator/internal/database"
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Rss was generated 2024-09-20 10:25:25 by https://xml-to-go.github.io/ in Ukraine.
type Rss struct {
	XMLName xml.Name `xml:"rss" json:"rss,omitempty"`
	Channel struct {
		Items []struct {
			Text        string `xml:",chardata" json:"text,omitempty"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Description string `xml:"description"`
		} `xml:"item" json:"item,omitempty"`
	} `xml:"channel" json:"channel,omitempty"`
}

func (cfg *apiConfig) fetchFeedContent(feed database.Feed, wg *sync.WaitGroup) {
	fmt.Printf("Requesting content from: %s (name: %s)\n", feed.Url, feed.Name)
	res, err := http.Get(feed.Url)
	if err != nil {
		fmt.Printf("error getting content from (name: %s) err: %v\n", feed.Name, err)
		wg.Done()
		return
	}
	rss := Rss{}
	defer res.Body.Close()
	fmt.Printf("Decoding content from:%s (name: %s)\n", feed.Url, feed.Name)
	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading content from (name: %s) err: %v\n", feed.Name, err)
		wg.Done()
		return
	}
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		fmt.Printf("error unmarshalling content from (name: %s) err: %v\n", feed.Name, err)
		wg.Done()
		return
	}
	fmt.Printf("Marking feed as fetched (name: %s)\n", feed.Name)
	cfg.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	})
	fmt.Printf("Marked (name: %s) feed as fetched\nFirst Post Title:%s\n", feed.Name, rss.Channel.Items[0].Title)
	fmt.Printf("Adding posts to db...\n")
	for _, post := range rss.Channel.Items {
		publishedAt := sql.NullTime{}
		if post.PubDate != "" {
			const timeLayout = "Mon, 2 Jan 2006 15:04:05 -0700"
			pubTime, err := time.Parse(timeLayout, post.PubDate)
			if err != nil {
				fmt.Printf("Error parsing published at time: %v\n", err)
			} else {
				publishedAt = sql.NullTime{
					Time:  pubTime,
					Valid: true,
				}
			}
		}
		description := sql.NullString{}
		if post.Description != "" {
			description = sql.NullString{
				String: post.Description,
				Valid:  true,
			}
		}
		err := cfg.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			PublishedAt: publishedAt,
			Title:       post.Title,
			Url:         post.Link,
			Description: description,
			FeedID:      feed.ID,
		})
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key value violates") {
				fmt.Printf("Error adding post to db: %v\n", err)
			}
			wg.Done()
			return
		}
	}
	wg.Done()
}

func (cfg *apiConfig) fetchFeedsWorker(maxConcurrency int) {
	for true {
		errCount := 0
		feeds, err := cfg.DB.GetNextFeedsToFetch(context.Background(), int32(maxConcurrency))
		fmt.Printf("Got feeds to fetch from DB:\n[\n")
		for _, feed := range feeds {
			fmt.Printf("\tName: %s; URL: %s\n", feed.Name, feed.Url)
		}
		fmt.Printf("]\n")
		if err != nil {
			fmt.Printf("Error getting feeds in fetch worker: %v (waiting %d seconds)\n", err, errCount*5)
			time.Sleep(time.Duration(errCount*5) * time.Second)
		} else {
			errCount = 0
			wg := &sync.WaitGroup{}
			for _, feed := range feeds {
				wg.Add(1)
				go cfg.fetchFeedContent(feed, wg)
			}
			if len(feeds) > 0 {
				fmt.Printf("Waiting for wg\n")
				wg.Wait()
				fmt.Printf("Finished Fetching posts results\n")
			}
			time.Sleep(30 * time.Second)
		}
	}
}
