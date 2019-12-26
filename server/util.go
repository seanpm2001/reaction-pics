package server

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
)

// appCacheString returns a cache string that can be used to bust browser/CDN caches
func appCacheString(logger *zap.SugaredLogger) string {
	appFile := staticPath + "app.js"
	info, err := os.Stat(appFile)
	if err != nil {
		logger.Error(err)
		rollbar.Error(rollbar.ERR, err)
		return ""
	}
	cacheString := strconv.FormatInt(info.ModTime().Unix(), 10)
	return cacheString
}

// logURL is a closure that logs (to stdout) the url and query of requests
func logURL(
	targetFunc func(http.ResponseWriter, *http.Request),
	logger *zap.SugaredLogger,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		logger.Info(url)
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

type handlerWithDeps func(http.ResponseWriter, *http.Request, handlerDeps)

// handlerDeps is a struct of handler dependencies
type handlerDeps struct {
	logger         *zap.SugaredLogger
	board          *tumblr.Board
	appCacheString string
}

// handlerGenerator returns a struct that can generate wrapped http handler functions
type handlerGenerator struct {
	newrelicApp *newrelic.Application
	logger      *zap.SugaredLogger
	deps        handlerDeps
}

// newHandlerGenerator returns a new handlerGenerator
func newHandlerGenerator(board *tumblr.Board, newrelicApp *newrelic.Application, logger *zap.SugaredLogger) handlerGenerator {
	deps := handlerDeps{
		logger:         logger,
		board:          board,
		appCacheString: appCacheString(logger),
	}
	return handlerGenerator{
		newrelicApp: newrelicApp,
		logger:      logger,
		deps:        deps,
	}
}

// newHandlerFunc returns a http handler function
func (g handlerGenerator) newHandler(pattern string, handlerFunc handlerWithDeps,
) (string, http.Handler) {
	f := func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r, g.deps)
	}
	return newrelic.WrapHandle(g.newrelicApp, pattern, http.HandlerFunc(logURL(f, g.logger)))
}
