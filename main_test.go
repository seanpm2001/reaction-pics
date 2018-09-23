package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNewRelicApp(t *testing.T) {
	setupEnv()
	app := getNewRelicApp()
	assert.NotNil(t, app)
}
