package server

import (
	"fmt"
	"github.com/albertyw/reaction-pics/tumblr"
	// Used for getting tumblr env vars
	_ "github.com/joho/godotenv/autoload"
	newrelic "github.com/newrelic/go-agent"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	dataURLPath = "/data.json"
	maxResults  = 20
)

var serverDir = filepath.Join(os.Getenv("ROOT_DIR"), "server")
var staticPath = fmt.Sprintf("%s/static", serverDir)
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

func rewriteFS(targetFunc func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		if path == "/" {
			path = "index.htm"
		}
		r.URL.Path = path
		targetFunc(w, r)
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
	numPosts := 0
	for _, post := range posts {
		postData := strings.ToLower(post.Title)
		if strings.Contains(postData, query) {
			selectedPosts = append(selectedPosts, post)
			numPosts++
		}
		if numPosts >= maxResults {
			break
		}
	}
	html := tumblr.PostsToJSON(selectedPosts)
	fmt.Fprintf(w, html)
}

// Run starts up the HTTP server
func Run(postChan <-chan tumblr.Post, newrelicApp newrelic.Application) {
	go loadPosts(postChan)
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	staticFS := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/", logURL(staticFS)))
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, dataURLPath, logURL(dataURLHandler)))
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/search", logURL(searchHandler)))
	http.ListenAndServe(address, nil)
}

func loadPosts(postChan <-chan tumblr.Post) {
	for p := range postChan {
		postsMutex.Lock()
		posts = append(posts, p)
		postsMutex.Unlock()
	}
	posts = *tumblr.SortPostsByLikes(&posts)
}
