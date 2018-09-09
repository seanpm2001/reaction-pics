package main

import (
	"testing"
)

func TestGetNewRelicApp(t *testing.T) {
	setupEnv()
	app := getNewRelicApp()
	if app == nil {
		t.Fail()
	}
}
