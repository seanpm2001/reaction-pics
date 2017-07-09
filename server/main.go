package server

import (
	"fmt"
	"github.com/albertyw/reaction-pics/tumblr"
	// Used for getting tumblr env vars
	_ "github.com/joho/godotenv/autoload"
	newrelic "github.com/newrelic/go-agent"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	dataURLPath = "/data.json"
	maxResults  = 20
)

var serverDir = filepath.Join(os.Getenv("ROOT_DIR"), "server")
var staticPath = fmt.Sprintf("%s/static/", serverDir)
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

// rewriteFS wraps a static file handler so to rewrite to the static directory
// and the root path is rewritten to index.htm
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

func postHandler(w http.ResponseWriter, r *http.Request) {
	pathStrings := strings.Split(r.URL.Path, "/")
	postIDString := pathStrings[2]
	postID, err := strconv.ParseInt(postIDString, 10, 64)
	if err != nil {
		fmt.Println("Cannot parse post id")
		http.NotFound(w, r)
		return
	}
	var post tumblr.Post
	for _, p := range posts {
		if p.ID == postID {
			post = p
		}
	}
	if post == (tumblr.Post{}) {
		fmt.Println("Cannot find post")
		http.NotFound(w, r)
		return
	}
	templateData, err := ioutil.ReadFile(staticPath + "post.htm")
	if err != nil {
		http.Error(w, "Cannot read post template", 500)
		return
	}
	htmlTemplate, err := template.New("page").Parse(string(templateData))
	if err != nil {
		http.Error(w, "Cannot parse post template", 500)
		return
	}
	data := struct {
		Post *tumblr.Post
	}{
		Post: &post,
	}
	err = htmlTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Cannot execute post template", 500)
		return
	}
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
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/post/", logURL(postHandler)))
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
