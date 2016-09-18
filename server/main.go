package server

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "asdf", r.URL.Path[1:])
}

func Run() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
