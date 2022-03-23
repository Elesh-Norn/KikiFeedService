package main

import (
  "os"
  "net/http"
  "html/template"
  "fmt"
)
var config Config

type context struct {
  Entries []entry
  Title string
}
func build(writer http.ResponseWriter, req *http.Request) {
  // as en endpoint until I figure out how to just execute a file in Go
  t := template.Must(template.ParseFiles("templates/layout.html"))
  feeds := getFeeds(config.Adresses)
  entries:= getSortedEntries(feeds)
  file, err := os.Create("static/index.html")
  if err != nil {
    fmt.Fprint(writer, "Building failed.")
    fmt.Println(err)
    return
  }
  t.Execute(file, context{Entries: entries, Title: config.Title})
  fmt.Fprint(writer, "Built!")
}

func reader(writer http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(writer, "index.html")
}

func main(){
  config = load_config()
  http.Handle("/", http.FileServer(http.Dir("./static")))
  http.HandleFunc("/build", build)
  http.ListenAndServe(":8090", nil)
}
