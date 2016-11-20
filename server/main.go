package server

import (
	"fmt"
	"github.com/flosch/pongo2"
	"net/http"
	"os"
)

func getTemplateDir() string {
	return os.Getenv("SERVER_TEMPLATES")
}

var indexPath = fmt.Sprintf("%s/index.htm", getTemplateDir())
var index = pongo2.Must(pongo2.FromFile(indexPath))

func handler(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request
	err := index.ExecuteWriter(pongo2.Context{"query": r.FormValue("query")}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Run starts up the HTTP server
func Run() {
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	http.HandleFunc("/", handler)
	http.ListenAndServe(address, nil)
}
