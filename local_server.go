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
  t := template.Must(template.ParseFiles("layout.html"))
  result, _ := GetFeed("https://emberger.xyz/index.xml")

  t.Execute(writer, result)

}

func main() {
  http.HandleFunc("/hello", hello)
  http.HandleFunc("/headers", headers)
  http.HandleFunc("/reader", reader)
  http.ListenAndServe(":8090", nil)
}
