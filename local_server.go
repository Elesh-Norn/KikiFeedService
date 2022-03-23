package main

import (
  "fmt"
  "net/http"
  "html/template"
)


func hello (writer http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(writer, "hello\n")
}

func headers(writer http.ResponseWriter, req *http.Request) {
  for name, headers := range req.Header {
    for _, h := range headers {
      fmt.Fprintf(writer, "%v: %v\n", name, h)
    }
  }
}

func reader(writer http.ResponseWriter, req *http.Request) {
  addresses := []string{
    "https://emberger.xyz/index.xml",
    //"https://martinfowler.com/feed.atom",
  }
  t := template.Must(template.ParseFiles("layout.html"))
  feeds := getFeeds(addresses)
  entries:= getSortedEntries(feeds)
  t.Execute(writer, entries)

}

func main() {
  http.HandleFunc("/hello", hello)
  http.HandleFunc("/headers", headers)
  http.HandleFunc("/reader", reader)
  http.ListenAndServe(":8090", nil)
}
