package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadPostsFromCSV(t *testing.T) {
	posts := ReadPostsFromCSV(getCSV(true))
	assert.Equal(t, len(posts), 1)
	assert.Equal(t, posts[0].ID, int64(1234))
	assert.Equal(t, posts[0].Title, "title")
	assert.Equal(t, posts[0].URL, "url")
	assert.Equal(t, posts[0].Image, "https://img.reaction.pics/file/reaction-pics/abcd.gif")
	assert.Equal(t, posts[0].Likes, int64(123))
}

func TestReadPercentFromCSV(t *testing.T) {
	data := []byte(`1234,a% b,url,image,123`)
	posts := ReadPostsFromCSV(data)
	assert.Equal(t, posts[0].Title, "a% b")
}
