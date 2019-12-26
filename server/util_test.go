package server

import (
	"testing"

	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCacheString(t *testing.T) {
	logger := zap.NewNop().Sugar()

	tempStaticPath := staticPath
	defer func() { staticPath = tempStaticPath }()
	cacheString := appCacheString(logger)
	assert.NotEqual(t, cacheString, "")

	staticPath = ""
	cacheString = appCacheString(logger)
	assert.Equal(t, cacheString, "")
}

func TestNewHandlerGenerator(t *testing.T) {
	b := tumblr.NewBoard([]tumblr.Post{})
	n := newrelic.Application{}
	l := zap.NewNop().Sugar()
	s := appCacheString(l)
	generator := newHandlerGenerator(&b, &n, l)
	assert.Equal(t, generator.newrelicApp, &n)
	assert.Equal(t, generator.logger, l)
	assert.Equal(t, generator.deps.logger, l)
	assert.Equal(t, generator.deps.board, &b)
	assert.Equal(t, generator.deps.appCacheString, s)
}
