package main 

import (
        "io"
        "log"
        "time"
        "github.com/mmcdole/gofeed"
        "sort"
        "sync"
        "regexp"
        "net/http"
)

type Feed = gofeed.Feed
type Item = gofeed.Item
type Parser = gofeed.Parser

type entry struct {
  BlogTitle string
  BlogLink string
  Title string
  Link string
  Published time.Time
  Description string
  }

type feedResponse struct {
  Url string
  LastModified string
  ETag string
  Body string

}

func createEntry(blogTitle string, blogLink string, item *Item) entry {
  // I just flatten the data struct to be able to have simpler template
  // and simpler time to sort them
  description := item.Description
  r := regexp.MustCompile("[^<]+[/>$]")
  if r.MatchString(description) {
    description = ""
  }
  // Sometimes rss feed put their whole article in html format in the
  // description field. I don't want to see the article in my tool.
  e := entry{
    BlogTitle: blogTitle,
    BlogLink: blogLink,
    Title: item.Title,
    Link: item.Link,
    Published: *item.PublishedParsed,
    Description: description,
  }
  return e
}

func getResponses(urls []string, lastResponses map[string]feedResponse, userAgent string) ([]feedResponse){
  // Get a list of Feeds
  result := []feedResponse{}
  var wg sync.WaitGroup
  
  for _, url := range(urls){
    
    wg.Add(1)
    // To avoid passing the same instance of the variable to 
    // each closure we need to initialise  an ew variable
    // https://go.dev/doc/faq#closures_and_goroutines
    url := url  
    go func() {
      defer wg.Done()
      resp, err := fetchFeed(url, lastResponses, userAgent)
      if err != nil {
        return
      }
      result = append(result, resp)
    }()
  }

  wg.Wait()
  
  return result
}


func fetchFeed(url string, lastResponses map[string]feedResponse, userAgent string) (feedResponse, error ){
  
  client := http.Client{}
  req, err := http.NewRequest("GET", url, nil)
  req.Header.Add("User-Agent", userAgent)
  oldResponse , ok := lastResponses[url]

  if ok {
    req.Header.Add("If-None-Match", oldResponse.ETag)
    req.Header.Add("If-Modified-Since", oldResponse.LastModified)
  }
  
  resp, err := client.Do(req)
  if err != nil {
    log.Printf("Error with "+ url)
    log.Printf(resp.Status, url)
    if ok {
      return oldResponse, nil
    }
    return feedResponse{"", "", "", ""}, err
  }
  if resp.StatusCode == 304 {
    log.Printf("Wasn't modified since last time " + url)
    return oldResponse, nil
  }
  bytes, _ := io.ReadAll(resp.Body)
  b := string(bytes)
  newLastResponse := feedResponse{url, resp.Header.Get("Last-Modified"), resp.Header.Get("Etag"),b}
  log.Printf("Succesfully fetched new "+ url)
  return newLastResponse, nil
}

func parseFeed(url string, body string, parser *Parser) (*Feed, error) {
  // Get 1 Feed
  feed, err := parser.ParseString(body)
  if err != nil {
    log.Printf("Something went wrong parsing the feed at requested url: %v \n", url)
    return nil , err
  }
  return feed, nil
}


func parseFeeds(lastResponses []feedResponse) []*Feed {
  // Get a list of Feeds
  result := make([]*Feed, 0)
  parser := gofeed.NewParser()
  for _, r := range(lastResponses){
      feed, err := parseFeed(r.Url, r.Body, parser)
      if err != nil {
        continue
      }
      result = append(result, feed)
    }
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
