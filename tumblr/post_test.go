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
	if post.id != 1234 {
		t.Fail()
	}
	if post.title != "title" {
		t.Fail()
	}
	if post.url != "url" {
		t.Fail()
	}
}
