package main 

import (
        "fmt"
        "time"
        "github.com/mmcdole/gofeed"
        "sort"
        "sync"
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
  // I just flatten the data struct to be able to have simpler template
  // and simpler time to sort them
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
    fmt.Printf("Something went wrong parsing the feed at requested url: %v \n", url)
    return nil , err
  }
  fmt.Printf("Succesfully fetched requested url: %v \n", url)
  return feed, nil
}

func getFeeds(urls []string) []*Feed {
  // Get a list of Feeds
  result := make([]*Feed, 0)
  
  var wg sync.WaitGroup
  
  for _, url := range(urls){
    
    wg.Add(1)
    // To avoid passing the same instance of the variable to 
    // each closure we need to initialise  an ew variable
    // https://go.dev/doc/faq#closures_and_goroutines
    url := url  
    go func() {
      defer wg.Done()
      feed, err := getFeed(url)
      if err != nil {
        return
      }
      result = append(result, feed)
    }()

  }

  wg.Wait()
  
  return result
}


func getSortedEntries(feeds []*Feed) []entry {
  // Put all the feeds into a big slice and transform them into entries
  result := make([]entry, 0)
  for _, feed := range feeds {
    result = append(result, getEntriesForFeed(config.ArticleNumber, feed)...)
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
