package server

import (
	"testing"

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
