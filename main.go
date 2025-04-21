package main

import (
  "log"
  "flag"
  "os"
  "gopkg.in/yaml.v2"
  "html/template"
  "path/filepath"
  "encoding/json"
)
type Config struct {
  Adresses []string
  Videos []string
  ArticleNumber int
  Title string
  UserAgent string
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

func load_last_Responses() map[string]feedResponse{
  store_path := filepath.Join(getExecutableDirPath(), "stored_responses.json")
  file, err := os.ReadFile(store_path)
  if err != nil {
    panic(err)
  }
  var responses map[string]feedResponse
  if err = json.Unmarshal(file, &responses); err != nil {
    panic(err)
  }
  return responses
}

var config Config
var lastResponses map[string]feedResponse
var newLastResponses []feedResponse
var server = flag.Bool("server", false, "Build and then launch local server")
var port = flag.Int("port", 8090, "Port for the server")
var help = flag.Bool("help", false, "Show help")

func build(templatePath string, staticPath string, topic []string) {
  // build the static site after polling each feed
  tPath := filepath.Join(getExecutableDirPath(), templatePath)
  sPath := filepath.Join(getExecutableDirPath(), staticPath )
  
  t := template.Must(template.ParseFiles(tPath))
  responses := getResponses(topic, lastResponses, config.UserAgent)
  for _, r := range(responses){
    newLastResponses = append(newLastResponses, r)
  }
  feeds := parseFeeds(responses)
  entries:= getSortedEntries(feeds)
  
  file, err := os.Create(sPath)
  logFile := getLogFile()
  log.SetOutput(logFile)

  if err != nil {
    log.Println("Text entries building failed.")
    log.Println(err)
    return
  }
  
  t.Execute(file, context{Entries: entries, Title: config.Title})

  log.Println("Built!")
}

func saveResponsesAsJson(responses []feedResponse){
  storePath := filepath.Join(getExecutableDirPath(), "stored_responses.json")
  responseMap := make(map[string]feedResponse)
  log.Println("Start")
  for _ , v := range(newLastResponses){
    log.Println(v.Url)
    responseMap[v.Url] = v
  }
  file, _ := os.OpenFile(storePath, os.O_CREATE|os.O_WRONLY, 0644)
  defer file.Close()
  encoder := json.NewEncoder(file)
  encoder.Encode(responseMap)
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
  lastResponses = load_last_Responses()
  build("templates/layout.html", "static/index.html", config.Adresses)
  build("templates/video_feed.html", "static/video_feed.html", config.Videos)
  saveResponsesAsJson(newLastResponses)
  if *help {
    flag.Usage()
    os.Exit(0)
  }

  if *server{
    launch_server(*port)
  }
}
