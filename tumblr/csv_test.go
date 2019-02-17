package tumblr

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	data := []byte("1234,title,url,abcd.gif,123")
	err = ioutil.WriteFile(getCSVPath(), data, 0644)
	assert.NoError(t, err)

	posts := ReadPostsFromCSV()
	assert.Equal(t, len(posts), 1)
	assert.Equal(t, posts[0].ID, int64(1234))
	assert.Equal(t, posts[0].Title, "title")
	assert.Equal(t, posts[0].URL, "url")
	assert.Equal(t, posts[0].Image, "http://static.reaction.pics/img/abcd.gif")
	assert.Equal(t, posts[0].Likes, int64(123))
}
