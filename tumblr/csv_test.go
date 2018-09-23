package tumblr

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const (
	testBlog = "testBlog"
)

var post = Post{
	1234,
	"title",
	"url",
	"http://placehold.it/350x150",
	123,
}

func cleanup() {
	csvLocation := getCSVPath()
	os.Remove(csvLocation)
}

func TestReadPostsFromCSV(t *testing.T) {
	defer cleanup()
	dotenvPath := os.Getenv("ROOT_DIR") + "/.env"
	err := godotenv.Load(dotenvPath)
	if err != nil {
		t.Fail()
	}
	data := []byte("1234,title,url,http://placehold.it/350x150,123")
	err = ioutil.WriteFile(getCSVPath(), data, 0644)
	if err != nil {
		t.Fail()
	}
	posts := ReadPostsFromCSV()
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
