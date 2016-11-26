package tumblr

import (
	"os"
	"strings"
	"testing"
)

var post = Post{
	1234,
	"title",
	"url",
}

func TestReadPostsFromCSV(t *testing.T) {
	posts := []Post{post}
	WritePostsToCSV(posts)
	posts = ReadPostsFromCSV()
	if len(posts) != 1 {
		t.Fail()
	}
	if posts[0].ID != 1234 {
		t.Fail()
	}
	if posts[0].Title != "title" {
		t.Fail()
	}
	if posts[0].URL != "url" {
		t.Fail()
	}
}

func TestWritePoststoCSV(t *testing.T) {
	posts := []Post{post}
	WritePostsToCSV(posts)
	file, err := os.Open(csvLocation)
	defer file.Close()
	if err != nil {
		t.Fail()
	}
	data := make([]byte, 20)
	length, err := file.Read(data)
	if length != 15 {
		t.Fail()
	}
	csvString := strings.TrimSpace(string(data[:14]))
	if csvString != "1234,title,url" {
		t.Fail()
	}

}

func TestGetRow(t *testing.T) {
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
