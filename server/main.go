package server

import (
	"fmt"
	"github.com/albertyw/reaction-pics/tumblr"
	// Used for getting tumblr env vars
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"os"
	"strings"
	"sync"
)

const dataURLPath = "/data.json"

var serverDir = os.Getenv("SERVER_DIR")
var indexPath = fmt.Sprintf("%s/templates/index.htm", serverDir)
var jsPath = fmt.Sprintf("%s/static/app.js", serverDir)
var cssPath = fmt.Sprintf("%s/static/global.css", serverDir)
var uRLFilePaths = map[string]func() (string, error){}
var posts []tumblr.Post
var postsMutex sync.RWMutex

// logURL is a closure that logs (to stdout) the url and query of requests
func logURL(
	targetFunc func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		fmt.Println(url)
		targetFunc(w, r)
	}
}

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
		file, err := os.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return
		}
		http.ServeContent(w, r, path, fileInfo.ModTime(), file)
	}
}

// dataURLHandler is an http handler for the dataURLPath response
func dataURLHandler(w http.ResponseWriter, r *http.Request) {
	postsMutex.RLock()
	defer postsMutex.RUnlock()
	html := tumblr.PostsToJSON(posts)
	fmt.Fprintf(w, html)
}

// searchHandler is an http handler to search data for keywords
// It matches the query against post titles and then ranks posts by number of likes
func searchHandler(w http.ResponseWriter, r *http.Request) {
	postsMutex.RLock()
	defer postsMutex.RUnlock()
	query := r.URL.Query().Get("query")
	query = strings.ToLower(query)
	selectedPosts := []tumblr.Post{}
	for _, post := range posts {
		postData := strings.ToLower(post.Title)
		if strings.Contains(postData, query) {
			selectedPosts = append(selectedPosts, post)
		}
	}
	selectedPosts = *tumblr.SortPostsByLikes(&selectedPosts)
	html := tumblr.PostsToJSON(selectedPosts)
	fmt.Fprintf(w, html)
}

// Run starts up the HTTP server
func Run(postChan chan tumblr.Post) {
	go loadPosts(postChan)
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	http.HandleFunc("/", logURL(exactURL(readFile(indexPath), "/")))
	http.HandleFunc("/static/app.js", logURL(readFile(jsPath)))
	http.HandleFunc("/static/global.css", logURL(readFile(cssPath)))
	http.HandleFunc(dataURLPath, logURL(dataURLHandler))
	http.HandleFunc("/search", logURL(searchHandler))
	http.ListenAndServe(address, nil)
}

func loadPosts(postChan chan tumblr.Post) {
	for p := range postChan {
		postsMutex.Lock()
		posts = append(posts, p)
		postsMutex.Unlock()
	}
}
