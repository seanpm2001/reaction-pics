package server

import (
	"fmt"
	"github.com/flosch/pongo2"
	"net/http"
)

const (
	templateDir = "server/templates/"
)

var indexPath = fmt.Sprintf("%s/index.htm", templateDir)
var index = pongo2.Must(pongo2.FromFile(indexPath))

func handler(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request
	err := index.ExecuteWriter(pongo2.Context{"query": r.FormValue("query")}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Run() {
	address := ":8080"
	fmt.Println("server listening on %s", address)
	http.HandleFunc("/", handler)
	http.ListenAndServe(address, nil)
}
