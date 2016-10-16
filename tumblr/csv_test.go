package tumblr

import (
	"testing"
)

func TestGetRow(t *testing.T) {
	post := Post{
		1234,
		"title",
		"url",
	}
	row := getRow(post)
	if row[0] != "1234" {
		t.Fail()
	}
	if row[1] != "title" {
		t.Fail()
	}
	if row[2] != "url" {
		t.Fail()
	}
}
