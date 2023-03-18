package server

import (
	"testing"

	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCacheString(t *testing.T) {
	cacheString := appCacheString()
	assert.NotEqual(t, cacheString, "")
}

func TestNewHandlerGenerator(t *testing.T) {
	b := tumblr.NewBoard([]tumblr.Post{})
	l := zap.NewNop()
	s := appCacheString()
	generator := newHandlerGenerator(&b, l)
	assert.Equal(t, generator.logger, l)
	assert.Equal(t, generator.deps.logger, l.Sugar())
	assert.Equal(t, generator.deps.board, &b)
	assert.Equal(t, generator.deps.appCacheString, s)
}
