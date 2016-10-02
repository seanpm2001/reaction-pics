package tumblr

import (
	_ "github.com/MariaTerzieva/gotumblr"
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
