package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	newrelic "github.com/newrelic/go-agent"
	"github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"
)

const (
	maxResults = 20
)

var serverDir = filepath.Join(os.Getenv("ROOT_DIR"), "server")
var staticPath = fmt.Sprintf("%s/static/", serverDir)
var uRLFilePaths = map[string]func() (string, error){}
var board = tumblr.NewBoard([]tumblr.Post{})

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
		r.URL.Path = path
		targetFunc(w, r)
	}
}

// indexHandler returns the index page template
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(staticPath + "index.htm")
	if err != nil {
		err = errors.New("Cannot read post template")
		fmt.Println(err)
		rollbar.RequestError(rollbar.ERR, r, err)
		http.Error(w, err.Error(), 500)
		return
	}
	templateData := struct {
		CacheString string
	}{appCacheString()}
	err = t.Execute(w, templateData)
	if err != nil {
		err = errors.Wrap(err, "Cannot execute template")
		fmt.Println(err)
		rollbar.RequestError(rollbar.ERR, r, err)
		http.Error(w, err.Error(), 500)
		return
	}
}

// searchHandler is an http handler to search data for keywords
// It matches the query against post titles and then ranks posts by number of likes
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	query = strings.ToLower(query)
	queriedBoard := board.FilterBoard(query)
	offsetString := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		offset = 0
	}
	data := map[string]interface{}{
		"offset":       offset,
		"totalResults": len(queriedBoard.Posts),
	}
	queriedBoard.LimitBoard(offset, maxResults)
	data["data"] = queriedBoard.PostsToJSON()
	dataBytes, _ := json.Marshal(data)
	fmt.Fprintf(w, string(dataBytes))
}

// postDataHandler is an http handler to return post data by ID
func postDataHandler(w http.ResponseWriter, r *http.Request) {
	pathStrings := strings.Split(r.URL.Path, "/")
	postIDString := pathStrings[2]
	postID, err := strconv.ParseInt(postIDString, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Cannot parse post id")
		fmt.Println(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}
	post := board.GetPostByID(postID)
	if post == nil {
		err = errors.New("Cannot find post")
		fmt.Println(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}
	data := map[string]interface{}{
		"offset":       0,
		"totalResults": 1,
		"data":         []tumblr.PostJSON{post.ToJSONStruct()},
	}
	marshalledPost, _ := json.Marshal(data)
	fmt.Fprintf(w, string(marshalledPost))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	pathStrings := strings.Split(r.URL.Path, "/")
	postIDString := pathStrings[2]
	postID, err := strconv.ParseInt(postIDString, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Cannot parse post id")
		fmt.Println(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}
	foundPost := false
	for _, p := range board.Posts {
		if p.ID == postID {
			foundPost = true
		}
	}
	if !foundPost {
		err = errors.New("Cannot find post")
		fmt.Println(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}
	indexHandler(w, r)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	postCount := strconv.Itoa(len(board.Posts))
	data := map[string]interface{}{
		"postCount": postCount,
		"keywords":  board.Keywords(),
	}
	stats, _ := json.Marshal(data)
	fmt.Fprintf(w, string(stats))
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	sm := stm.NewSitemap(0)
	sm.SetDefaultHost(os.Getenv("HOST"))

	sm.Create()
	sm.Add(stm.URL{{"loc", "/"}})
	for _, url := range board.URLs() {
		sm.Add(stm.URL{{"loc", url}})
	}
	w.Write(sm.XMLContent())
}

// Run starts up the HTTP server
func Run(postChan <-chan tumblr.Post, newrelicApp newrelic.Application) {
	go loadPosts(postChan)
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	staticFS := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/", logURL(indexHandler)))
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/static/", logURL(staticFS)))
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/search", logURL(searchHandler)))
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/postdata/", logURL(postDataHandler)))
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/post/", logURL(postHandler)))
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/stats.json", logURL(statsHandler)))
	http.HandleFunc(newrelic.WrapHandleFunc(newrelicApp, "/sitemap.xml", logURL(sitemapHandler)))
	http.ListenAndServe(address, nil)
}

func loadPosts(postChan <-chan tumblr.Post) {
	for p := range postChan {
		board.AddPost(p)
		board.SortPostsByLikes()
	}
}

func appCacheString() string {
	appFile := staticPath + "app.js"
	info, err := os.Stat(appFile)
	if err != nil {
		fmt.Println(err)
		rollbar.Error(rollbar.ERR, err)
		return ""
	}
	cacheString := strconv.FormatInt(info.ModTime().Unix(), 10)
	return cacheString
}
