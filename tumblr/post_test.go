package tumblr

import (
	"github.com/MariaTerzieva/gotumblr"
	"testing"
)

func TestPost(t *testing.T) {
	post := Post{
		1234,
		"title",
		"url",
	}
	if post.Id != 1234 {
		t.Fail()
	}
	if post.Title != "title" {
		t.Fail()
	}
	if post.Url != "url" {
		t.Fail()
	}
}

func TestTumblrToPost(t *testing.T) {
	tumblrPost := gotumblr.TextPost{
		Title: "title",
	}
	tumblrPost.Id = 1234
	tumblrPost.Post_url = "url"
	post := TumblrToPost(&tumblrPost)
	if post.Id != 1234 {
		t.Fail()
	}
	if post.Title != "title" {
		t.Fail()
	}
	if post.Url != "url" {
		t.Fail()
	}
}

func TestCSVToPost(t *testing.T) {
	row := make([]string, 3)
	row[0] = "1234"
	row[1] = "title"
	row[2] = "url"
	post := CSVToPost(row)
	if post.Id != 1234 {
		t.Fail()
	}
	if post.Title != "title" {
		t.Fail()
	}
	if post.Url != "url" {
		t.Fail()
	}
}
