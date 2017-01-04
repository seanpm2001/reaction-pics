package server

import (
	"errors"
	"fmt"
	"github.com/albertyw/devops-reactions-index/tumblr"
	"io/ioutil"
	"net/http"
	"os"
)

const dataURLPath = "/data.json"

var templateDir = os.Getenv("SERVER_TEMPLATES")
var indexPath = fmt.Sprintf("%s/index.htm", templateDir)
var uRLFilePaths = map[string]func() (string, error){}
var posts []tumblr.Post

func readFile(p string) func(http.ResponseWriter, *http.Request) {
	path := p
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return
		}
		html := string(data)
		fmt.Fprintf(w, html)
	}
}

func fakeHandler() (string, error) {
	return "", errors.New("Called Fake Handler")
}

func getURLHandler(urlPath string) (func() (string, error), error) {
	handler, exists := uRLFilePaths[urlPath]
	if !exists {
		return fakeHandler, errors.New("")
	}
	return handler, nil
}

func dataURLHandler(w http.ResponseWriter, r *http.Request) {
	html := tumblr.PostsToJSON(posts)
	fmt.Fprintf(w, html)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request
	urlPath := r.URL.Path
	handler, err := getURLHandler(urlPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	data, err := handler()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, data)
}

// Run starts up the HTTP server
func Run(p []tumblr.Post) {
	posts = p
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	http.HandleFunc("/", readFile(indexPath))
	http.HandleFunc(dataURLPath, dataURLHandler)
	http.ListenAndServe(address, nil)
}
