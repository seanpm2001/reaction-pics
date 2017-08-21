package main

import (
	"os"
	"testing"
)

func TestGetReadPostsFromTumblr(t *testing.T) {
	defer os.Unsetenv(readPostsFromTumblrEnv)
	os.Setenv(readPostsFromTumblrEnv, "false")
	if getReadPostsFromTumblr() {
		t.Fail()
	}
	os.Setenv(readPostsFromTumblrEnv, "True")
	if !getReadPostsFromTumblr() {
		t.Fail()
	}
}

func TestGetNewRelicApp(t *testing.T) {
	app := getNewRelicApp()
	if app == nil {
		t.Fail()
	}
}
