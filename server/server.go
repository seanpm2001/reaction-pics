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
	"go.uber.org/zap"
)

const (
	maxResults = 20
)

var serverDir = filepath.Join(os.Getenv("ROOT_DIR"), "server")
var staticPath = fmt.Sprintf("%s/static/", serverDir)
var board *tumblr.Board

// indexHandler is an http handler that returns the index page HTML
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && !strings.HasPrefix(r.URL.Path, "/post/") {
		err := fmt.Errorf("file not found: %s", r.URL.Path)
		fmt.Println(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}
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

// searchHandler is an http handler to search data for keywords in json format
// It matches the query against post titles and then ranks posts by number of likes
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	query = strings.ToLower(query)
	queriedBoard := board.FilterBoard(query)
	if query == "" {
		queriedBoard.RandomizePosts()
	}
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
	queriedBoard.SortPostsByLikes()
	data["data"] = queriedBoard.PostsToJSON()
	dataBytes, _ := json.Marshal(data)
	fmt.Fprint(w, string(dataBytes))
}

// postDataHandler is an http handler to return post data by ID in json format
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
	fmt.Fprint(w, string(marshalledPost))
}

// postHandler is an http handler that validates the correctness of a post url
// and returns the index page html to render it correct
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

// statsHandler returns internal stats about the reaction.pics DB as json
func statsHandler(w http.ResponseWriter, r *http.Request) {
	postCount := strconv.Itoa(len(board.Posts))
	data := map[string]interface{}{
		"postCount": postCount,
		"keywords":  board.Keywords(),
	}
	stats, _ := json.Marshal(data)
	fmt.Fprint(w, string(stats))
}

// sitemapHandler returns a sitemap of reaction.pics as an xml file
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

// staticHandler returns static files
func staticHandler(w http.ResponseWriter, r *http.Request) {
	staticFS := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	staticFS(w, r)
}

// Run starts up the HTTP server
func Run(newrelicApp newrelic.Application, logger *zap.SugaredLogger) {
	board = tumblr.InitializeBoard()
	address := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logger.Infof("server listening on %s", address)
	generator := newHandlerGenerator(newrelicApp, logger)
	http.Handle(generator.newHandler("/", indexHandler))
	http.Handle(generator.newHandler("/search", searchHandler))
	http.Handle(generator.newHandler("/postdata/", postDataHandler))
	http.Handle(generator.newHandler("/post/", postHandler))
	http.Handle(generator.newHandler("/stats.json", statsHandler))
	http.Handle(generator.newHandler("/sitemap.xml", sitemapHandler))
	http.Handle(generator.newHandler("/static/", staticHandler))
	http.ListenAndServe(address, nil)
}
