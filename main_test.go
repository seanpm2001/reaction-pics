package main

import (
	"os"
	"testing"

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
