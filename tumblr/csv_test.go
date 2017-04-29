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
	"http://placehold.it/350x150",
}

func cleanup() {
	os.Remove(csvLocation)
}

func TestReadPostsFromCSV(t *testing.T) {
	defer cleanup()
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
	if posts[0].Image != "http://placehold.it/350x150" {
		t.Fail()
	}
}

func TestWritePostsToCSV(t *testing.T) {
	defer cleanup()
	posts := []Post{post}
	WritePostsToCSV(posts)
	file, err := os.Open(csvLocation)
	defer file.Close()
	if err != nil {
		t.Fail()
	}
	data := make([]byte, 50)
	length, err := file.Read(data)
	if length != 43 {
		t.Fail()
	}
	csvString := strings.TrimSpace(string(data[:43]))
	if csvString != "1234,title,url,http://placehold.it/350x150" {
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
	if row[3] != "http://placehold.it/350x150" {
		t.Fail()
	}
}
