package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
)

const (
	maxResults = 20
)

type metaHeader struct {
	Property string
	Content  string
}

func relToAbsPath(path string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		filename = "."
	}
	absPath := filepath.Join(filepath.Dir(filename), path)
	return absPath
}

// indexHandler is an http handler that returns the index page HTML
func indexHandler(w http.ResponseWriter, r *http.Request, d handlerDeps) {
	if r.URL.Path != "/" && !strings.HasPrefix(r.URL.Path, "/post/") {
		err := fmt.Errorf("file not found: %s", r.URL.Path)
		d.logger.Warn(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}
	indexHandlerWithHeaders(w, r, d, []metaHeader{})
}

func indexHandlerWithHeaders(w http.ResponseWriter, r *http.Request, d handlerDeps, headers []metaHeader) {
	path := relToAbsPath("static/index.htm")
	t, err := template.ParseFiles(path)
	if err != nil {
		err = errors.Wrap(err, "Cannot read post template")
		d.logger.Error(err)
		rollbar.RequestError(rollbar.ERR, r, err)
		http.Error(w, err.Error(), 500)
		return
	}
	templateData := struct {
		CacheString string
		MetaHeaders []metaHeader
	}{
		CacheString: d.appCacheString,
		MetaHeaders: headers,
	}
	err = t.Execute(w, templateData)
	if err != nil {
		err = errors.Wrap(err, "Cannot execute template")
		d.logger.Error(err)
		rollbar.RequestError(rollbar.ERR, r, err)
		http.Error(w, err.Error(), 500)
		return
	}
}

// searchHandler is an http handler to search data for keywords in json format
// It matches the query against post titles and then ranks posts by number of likes
func searchHandler(w http.ResponseWriter, r *http.Request, d handlerDeps) {
	query := r.URL.Query().Get("query")
	query = strings.ToLower(query)
	queriedBoard := d.board.FilterBoard(query)
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
func postDataHandler(w http.ResponseWriter, r *http.Request, d handlerDeps) {
	pathStrings := strings.Split(r.URL.Path, "/")
	postIDString := pathStrings[2]
	postID, err := strconv.ParseInt(postIDString, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Cannot parse post id")
		d.logger.Warn(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}
	post := d.board.GetPostByID(postID)
	if post == nil {
		err = errors.New("Cannot find post")
		d.logger.Warn(err)
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
func postHandler(w http.ResponseWriter, r *http.Request, d handlerDeps) {
	pathStrings := strings.Split(r.URL.Path, "/")
	postIDString := pathStrings[2]
	postID, err := strconv.ParseInt(postIDString, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Cannot parse post id")
		d.logger.Warn(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}
	var post *tumblr.Post
	for _, p := range d.board.Posts {
		if p.ID == postID {
			post = &p
			break
		}
	}
	if post == nil {
		err = errors.New("Cannot find post")
		d.logger.Warn(err)
		rollbar.RequestError(rollbar.WARN, r, err)
		http.NotFound(w, r)
		return
	}

	headers := []metaHeader{
		metaHeader{"og:title", post.Title},
		metaHeader{"og:image", post.Image},
	}

	indexHandlerWithHeaders(w, r, d, headers)
}

// statsHandler returns internal stats about the reaction.pics DB as json
func statsHandler(w http.ResponseWriter, r *http.Request, d handlerDeps) {
	postCount := strconv.Itoa(len(d.board.Posts))
	data := map[string]interface{}{
		"postCount": postCount,
		"keywords":  d.board.Keywords(),
	}
	stats, _ := json.Marshal(data)
	fmt.Fprint(w, string(stats))
}

// sitemapHandler returns a sitemap of reaction.pics as an xml file
func sitemapHandler(w http.ResponseWriter, r *http.Request, d handlerDeps) {
	sm := stm.NewSitemap(0)
	sm.SetDefaultHost(os.Getenv("HOST"))

	sm.Create()
	sm.Add(stm.URL{{"loc", "/"}})
	for _, url := range d.board.URLs() {
		sm.Add(stm.URL{{"loc", url}})
	}
	w.Write(sm.XMLContent())
}

// staticHandler returns static files
func staticHandler(w http.ResponseWriter, r *http.Request, _ handlerDeps) {
	staticFS := rewriteFS(http.FileServer(http.Dir(relToAbsPath("static"))).ServeHTTP)
	staticFS(w, r)
}

func timeHandler(w http.ResponseWriter, r *http.Request, _ handlerDeps) {
	unixTime := int32(time.Now().Unix())
	data := map[string]interface{}{
		"unixtime": unixTime,
	}
	timeData, _ := json.Marshal(data)
	fmt.Fprint(w, string(timeData))
}

func faviconHandler(w http.ResponseWriter, r *http.Request, _ handlerDeps) {
	faviconPath := relToAbsPath("static/favicon/favicon.iso")
	http.ServeFile(w, r, faviconPath)
}

func robotsTxtHandler(w http.ResponseWriter, r *http.Request, _ handlerDeps) {
	fmt.Fprint(w, "")
}

// Run starts up the HTTP server
func Run(newrelicApp *newrelic.Application, logger *zap.SugaredLogger) {
	board := tumblr.InitializeBoard()
	address := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logger.Infof("server listening on %s", address)
	generator := newHandlerGenerator(board, newrelicApp, logger)
	http.Handle(generator.newHandler("/", indexHandler))
	http.Handle(generator.newHandler("/favicon.ico", faviconHandler))
	http.Handle(generator.newHandler("/robots.txt", robotsTxtHandler))
	http.Handle(generator.newHandler("/search", searchHandler))
	http.Handle(generator.newHandler("/postdata/", postDataHandler))
	http.Handle(generator.newHandler("/post/", postHandler))
	http.Handle(generator.newHandler("/stats.json", statsHandler))
	http.Handle(generator.newHandler("/sitemap.xml", sitemapHandler))
	http.Handle(generator.newHandler("/static/", staticHandler))
	http.Handle(generator.newHandler("/time/", timeHandler))
	http.ListenAndServe(address, nil)
}
