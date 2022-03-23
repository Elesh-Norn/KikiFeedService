package main

import (
        "fmt"
        "github.com/mmcdole/gofeed"
        "strings"
)

type Feed = gofeed.Feed

func GetFeed(url string) (*Feed, error) {
  parser := gofeed.NewParser()
  feed, err := parser.ParseURL(url)
  
  if err != nil {
    fmt.Println("Something went wrong parsing the feed at requested url")
    return nil , err
  }
  return feed, nil
}

func OutputFeed() string {
  feed, err := GetFeed("https://emberger.xyz/index.xml")
  if err != nil {
    panic(err)
  }
  max_lenght := 10
  if feed.Len() < max_lenght {
    max_lenght = feed.Len()
  }

  result := ""
  
  for i:= 0; i <= max_lenght - 1; i++ {
    article := feed.Items[i]
    result = strings.Join([]string{ result, "------------", article.Title, article.Description}, "\n")
  }
  if result == "" {
    result = "Nothing found."
  }
  return result
}
