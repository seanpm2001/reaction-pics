package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheString(t *testing.T) {
	tempStaticPath := staticPath
	defer func() { staticPath = tempStaticPath }()
	cacheString := appCacheString()
	assert.NotEqual(t, cacheString, "")

	staticPath = ""
	cacheString = appCacheString()
	assert.Equal(t, cacheString, "")
}
