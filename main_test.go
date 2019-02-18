package main

import (
	"os"
	"testing"

	"github.com/rollbar/rollbar-go"
	"github.com/stretchr/testify/assert"
)

func TestSetupEnv(t *testing.T) {
	setupEnv()
	port := os.Getenv("PORT")
	assert.NotEqual(t, port, "")
}

func TestGetNewRelicApp(t *testing.T) {
	setupEnv()
	app := getNewRelicApp()
	assert.NotNil(t, app)
}

func TestSetupRollbar(t *testing.T) {
	setupRollbar()
	assert.Equal(t, rollbar.Token(), os.Getenv("ROLLBAR_SERVER_TOKEN"))
	assert.Equal(t, rollbar.Environment(), os.Getenv("ENVIRONMENT"))
}
