package tumblr

import (
	"io/ioutil"
	"os"
	"testing"
)

const testBlog = "testBlog"

var post = Post{
	1234,
	"title",
	"url",
	"http://placehold.it/350x150",
	123,
}

func cleanup() {
	csvLocation := getCSVPath(testBlog)
	os.Remove(csvLocation)
}

func TestReadPostsFromCSV(t *testing.T) {
	defer cleanup()
	data := []byte("1234,title,url,http://placehold.it/350x150,123")
	ioutil.WriteFile(getCSVPath(testBlog), data, 0644)
	posts := ReadPostsFromCSV(testBlog)
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
	if posts[0].Likes != 123 {
		t.Fail()
	}
}
