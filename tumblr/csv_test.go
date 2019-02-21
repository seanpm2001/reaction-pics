package tumblr

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Rename(getCSVPath()+".bak", getCSVPath())
}

func TestReadPostsFromCSV(t *testing.T) {
	os.Rename(getCSVPath(), getCSVPath()+".bak")
	defer cleanup()

	data := []byte("1234,title,url,abcd.gif,123")
	err := ioutil.WriteFile(getCSVPath(), data, 0644)
	assert.NoError(t, err)

	posts := ReadPostsFromCSV()
	assert.Equal(t, len(posts), 1)
	assert.Equal(t, posts[0].ID, int64(1234))
	assert.Equal(t, posts[0].Title, "title")
	assert.Equal(t, posts[0].URL, "url")
	assert.Equal(t, posts[0].Image, "http://static.reaction.pics/img/abcd.gif")
	assert.Equal(t, posts[0].Likes, int64(123))
}
