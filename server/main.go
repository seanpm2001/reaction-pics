package server

import (
	"fmt"
	"github.com/albertyw/devops-reactions-index/tumblr"
	// Used for getting tumblr env vars
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const dataURLPath = "/data.json"

var templateDir = os.Getenv("SERVER_TEMPLATES")
var indexPath = fmt.Sprintf("%s/index.htm", templateDir)
var uRLFilePaths = map[string]func() (string, error){}
var posts []tumblr.Post

// exactURL is a closure that checks that the http match is an exact url path
// match instead of allowing for net/http's loose match
func exactURL(
	targetFunc func(http.ResponseWriter, *http.Request),
	requestedPath string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != requestedPath {
			http.NotFound(w, r)
			return
		}
		targetFunc(w, r)
		return
	}
}

// readFile returns a function that reads the file at a given path and makes a
// response from it
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

// dataURLHandler is an http handler for the dataURLPath response
func dataURLHandler(w http.ResponseWriter, r *http.Request) {
	html := tumblr.PostsToJSON(posts)
	fmt.Fprintf(w, html)
}

// searchHandler is an http handler to search data for keywords
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	query = strings.ToLower(query)
	selectedPosts := []tumblr.Post{}
	for _, post := range posts {
		postData := strings.ToLower(post.Title)
		if strings.Contains(postData, query) {
			selectedPosts = append(selectedPosts, post)
		}
	}
	html := tumblr.PostsToJSON(selectedPosts)
	fmt.Fprintf(w, html)
}

// Run starts up the HTTP server
func Run(p []tumblr.Post) {
	posts = p
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	http.HandleFunc("/", exactURL(readFile(indexPath), "/"))
	http.HandleFunc(dataURLPath, dataURLHandler)
	http.HandleFunc("/search", searchHandler)
	http.ListenAndServe(address, nil)
}
