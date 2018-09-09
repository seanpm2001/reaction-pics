package main

import (
	"testing"
)

func TestGetNewRelicApp(t *testing.T) {
	app := getNewRelicApp()
	if app == nil {
		t.Fail()
	}
}
