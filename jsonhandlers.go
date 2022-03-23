package main

import (
  "encoding/json"
  "os"
)

func dump_entry(entries []entry) []byte{
  serial, _ := json.Marshal(entries)
  return serial
}

func dump_entries(){
  addresses := []string{
    "https://emberger.xyz/index.xml",
    //"https://martinfowler.com/feed.atom",
  }
  feeds := getFeeds(addresses)
  entries:= getSortedEntries(feeds)
  json := dump_entry(entries)
  err := os.WriteFile("/tmp/feed_dump.json", json, 0644)
  if err != nil {
    panic(err)
  }
}

func load_entries() []entry {
  file, err := os.ReadFile("/tmp/feed_dump.json")
  if err != nil {
    panic(err)
  }

  var entries []entry
  if err := json.Unmarshal(file, &entries);
  err != nil {
    panic(err)
  }
  return entries
}

