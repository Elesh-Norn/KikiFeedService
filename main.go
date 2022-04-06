package main

import (
  "log"
  "flag"
  "os"
  "gopkg.in/yaml.v2"
  "html/template"
  "path/filepath"
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
  configPath := filepath.Join(getExecutableDirPath(), "config.yaml")
  file, err := os.ReadFile(configPath)
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
  tPath := filepath.Join(getExecutableDirPath(), "templates/layout.html")
  sPath := filepath.Join(getExecutableDirPath(), "static/index.html")
  t := template.Must(template.ParseFiles(tPath))
  feeds := getFeeds(config.Adresses)
  entries:= getSortedEntries(feeds)
  file, err := os.Create(sPath)
  logFile := getLogFile()
  log.SetOutput(logFile)

  if err != nil {
    log.Println("Building failed.")
    log.Println(err)
    return
  }
  
  t.Execute(file, context{Entries: entries, Title: config.Title})
  log.Println("Built!")
}

func getExecutableDirPath() (string) {
  ex, err := os.Executable()
  if err != nil {
    panic(err)
  }
  return filepath.Dir(ex)
}

func getLogFile() (*os.File) {
  logPath := filepath.Join(getExecutableDirPath(), "KikiFeedService.log")
  file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err !=nil {
    log.Fatal(err)
  }
  return file
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
