package tumblr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadPostsFromCSV(t *testing.T) {
	posts := ReadPostsFromCSV(getCSVPath(true))
	assert.Equal(t, len(posts), 1)
	assert.Equal(t, posts[0].ID, int64(1234))
	assert.Equal(t, posts[0].Title, "title")
	assert.Equal(t, posts[0].URL, "url")
	assert.Equal(t, posts[0].Image, "http://static.reaction.pics/img/abcd.gif")
	assert.Equal(t, posts[0].Likes, int64(123))
}
