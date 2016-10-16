package server

import (
	"fmt"
	"github.com/flosch/pongo2"
	"net/http"
	"os"
)

func getTemplateDir() string {
	goPath := os.Getenv("GOPATH")
	currentDir := fmt.Sprintf("%s/src/github.com/albertyw/devops-reactions-index", goPath)
	return fmt.Sprintf("%s/server/templates", currentDir)
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
	address := ":8080"
	fmt.Println("server listening on", address)
	http.HandleFunc("/", handler)
	http.ListenAndServe(address, nil)
}
