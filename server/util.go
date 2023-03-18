package server

import (
	"crypto/md5"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/albertyw/reaction-pics/tumblr"
	"go.uber.org/zap"
)

// appCacheString returns a cache string that can be used to bust browser/CDN caches
func appCacheString() string {
	appJS, err := staticFiles.ReadFile("static/app.js")
	if err != nil {
		return "notfound"
	}
	cacheHash := md5.Sum(appJS)
	return fmt.Sprintf("%x", cacheHash)
}

// logURL is a closure that logs (to stdout) the url and query of requests
func logURL(
	targetFunc http.HandlerFunc,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		logger.Info(url)
		targetFunc(w, r)
	}
}

type handlerWithDeps func(http.ResponseWriter, *http.Request, handlerDeps)

// handlerDeps is a struct of handler dependencies
type handlerDeps struct {
	logger         *zap.Logger
	board          *tumblr.Board
	appCacheString string
}

// handlerGenerator returns a struct that can generate wrapped http handler functions
type handlerGenerator struct {
	logger *zap.Logger
	deps   handlerDeps
}

// newHandlerGenerator returns a new handlerGenerator
func newHandlerGenerator(board *tumblr.Board, logger *zap.Logger) handlerGenerator {
	deps := handlerDeps{
		logger:         logger,
		board:          board,
		appCacheString: appCacheString(),
	}
	return handlerGenerator{
		logger: logger,
		deps:   deps,
	}
}

// newHandlerFunc returns a http handler function
func (g handlerGenerator) newHandler(handlerFunc handlerWithDeps) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r, g.deps)
	}
	return http.HandlerFunc(logURL(f, g.logger))
}
