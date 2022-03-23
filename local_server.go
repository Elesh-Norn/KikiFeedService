package main 

import (
  "net/http"
  "fmt"
  "strconv"
)

func reader(writer http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(writer, "index.html")
}

func launch_server(port int) {
  fmt.Println("Launching local server on " + strconv.Itoa(port))
  http.Handle("/", http.FileServer(http.Dir("./static")))
  http.ListenAndServe(":" + strconv.Itoa(port), nil)
}
