package main

import (
  "net/http"
  "html/template"
)

func reader(writer http.ResponseWriter, req *http.Request) {
  t := template.Must(template.ParseFiles("layout.html"))
  entries:= load_entries()
  t.Execute(writer, entries)
}

func main(){
  http.HandleFunc("/", reader)
  http.ListenAndServe(":8090", nil)
}
