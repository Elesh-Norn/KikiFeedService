package main

import (
  "fmt"
  "flag"
  "os"
  "gopkg.in/yaml.v2"
  "html/template"
)
type Config struct {
  Adresses []string
  ArticleNumber int
  Title string
}

type context struct {
  Entries []entry
  Title string
}

func load_config() Config{
  file, err := os.ReadFile("config.yaml")
  if err != nil {
    panic(err)
  }
  var c Config
  if err := yaml.Unmarshal(file, &c);
  err != nil {
    panic(err)
  }
  return c
}

var config Config
var server = flag.Bool("server", false, "Build and then launch local server")
var port = flag.Int("port", 8090, "Port for the server")
var help = flag.Bool("help", false, "Show help")

func build() {
  // build the static site after polling each feed
  t := template.Must(template.ParseFiles("templates/layout.html"))
  feeds := getFeeds(config.Adresses)
  entries:= getSortedEntries(feeds)
  file, err := os.Create("static/index.html")
  
  if err != nil {
    fmt.Println("Building failed.")
    fmt.Println(err)
    return
  }
  
  t.Execute(file, context{Entries: entries, Title: config.Title})
  fmt.Println("Built!")
}

func main() {
  flag.Parse()
  
  config = load_config()
  build()
  if *help {
    flag.Usage()
    os.Exit(0)
  }

  if *server{
    launch_server(*port)
  }
}
