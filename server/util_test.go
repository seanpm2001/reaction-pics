package server

import (
	"testing"

	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCacheString(t *testing.T) {
	logger := zap.NewNop().Sugar()

	cacheString := appCacheString(logger)
	assert.NotEqual(t, cacheString, "")
}

func TestNewHandlerGenerator(t *testing.T) {
	b := tumblr.NewBoard([]tumblr.Post{})
	l := zap.NewNop().Sugar()
	s := appCacheString(l)
	generator := newHandlerGenerator(&b, l)
	assert.Equal(t, generator.logger, l)
	assert.Equal(t, generator.deps.logger, l)
	assert.Equal(t, generator.deps.board, &b)
	assert.Equal(t, generator.deps.appCacheString, s)
}
