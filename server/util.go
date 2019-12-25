package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	newrelic "github.com/newrelic/go-agent"
	"github.com/rollbar/rollbar-go"
)

// appCacheString returns a cache string that can be used to bust browser/CDN caches
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

// handlerGenerator returns a struct that can generate wrapped http handler functions
type handlerGenerator struct {
	newrelicApp newrelic.Application
}

// newHandlerGenerator returns a new handlerGenerator
func newHandlerGenerator(newrelicApp newrelic.Application) handlerGenerator {
	return handlerGenerator{
		newrelicApp: newrelicApp,
	}
}

// newHandlerFunc returns a http handler function
func (g handlerGenerator) newHandlerFunc(pattern string, handlerFunc func(http.ResponseWriter, *http.Request),
) (string, func(http.ResponseWriter, *http.Request)) {
	return newrelic.WrapHandleFunc(g.newrelicApp, pattern, logURL(handlerFunc))
}
