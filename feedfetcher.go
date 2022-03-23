package main

import (
        "fmt"
        "time"
        "github.com/mmcdole/gofeed"
        "sort"
)

type Feed = gofeed.Feed
type Item = gofeed.Item

type entry struct {
  BlogTitle string
  BlogLink string
  Title string
  Link string
  Published time.Time
  Description string
}

func createEntry(blogTitle string, blogLink string, item *Item) entry {
  e := entry{
    BlogTitle: blogTitle,
    BlogLink: blogLink,
    Title: item.Title,
    Link: item.Link,
    Published: *item.PublishedParsed,
    Description: item.Description,
  }
  return e
}


func getFeed(url string) (*Feed, error) {
  // Get 1 Feed
  parser := gofeed.NewParser()
  feed, err := parser.ParseURL(url)
  
  if err != nil {
    fmt.Println("Something went wrong parsing the feed at requested url: %v", url)
    return nil , err
  }
  return feed, nil
}

func getFeeds(urls []string) []*Feed {
  // Get a list of Feeds
  result := make([]*Feed, 0)
  for _, url := range(urls){
    feed, err := getFeed(url)
    if err != nil {
      // Skipping faulty feed
      continue
    }
    result = append(result, feed)
  }
  return result
}


func getSortedEntries(feeds []*Feed) []entry {
  // Put all the feeds into a big slice and transform them into entrie
  result := make([]entry, 0)
  for _, feed := range feeds {
    result = append(result, getEntriesForFeed(10, feed)...)
  }

  // Sort the entries from most recent to most ancient
  sort.Slice(result, func(i, j int) bool{
    return result[i].Published.After(result[j].Published)
  })
  return result
}

func getEntriesForFeed(max int, feed *Feed) []entry {
  result := make([]entry, 0)
  if feed.Len() < max {
    max = feed.Len()
  }
  for _, e := range feed.Items[:max] {
    result = append(result, createEntry(feed.Title, feed.Link, e))
  }
  return result
}
