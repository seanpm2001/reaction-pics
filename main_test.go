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

func TestSetupRollbar(t *testing.T) {
	setupRollbar()
	assert.Equal(t, rollbar.Token(), os.Getenv("ROLLBAR_SERVER_TOKEN"))
	assert.Equal(t, rollbar.Environment(), os.Getenv("ENVIRONMENT"))
}

func TestGetLogger(t *testing.T) {
	origEnv := os.Getenv("ENVIRONMENT")
	defer func() { os.Setenv("ENVIRONMENt", origEnv) }()

	os.Setenv("ENVIRONMENT", "development")
	logger := getLogger()
	assert.NotNil(t, logger)

	os.Setenv("ENVIRONMENT", "production")
	logger = getLogger()
	assert.NotNil(t, logger)
}
