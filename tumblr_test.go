package main

import (
	"github.com/MariaTerzieva/gotumblr"
	"testing"
)

var mockBlogInfo = gotumblr.BlogInfo{
	Title:       "title",
	Posts:       10,
	Name:        "name",
	Url:         "url",
	Updated:     5,
	Description: "description",
	Ask:         true,
	Ask_anon:    false,
	Likes:       50,
}

func TestParsePosts(t *testing.T) {

}
