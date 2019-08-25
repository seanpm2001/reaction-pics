package tumblr

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	post := Post{
		1234,
		"title",
		"url",
		"http://static.reaction.pics/img/abcd.gif",
		123,
	}
	assert.Equal(t, post.ID, int64(1234))
	assert.Equal(t, post.Title, "title")
	assert.Equal(t, post.URL, "url")
	assert.Equal(t, post.Image, "http://static.reaction.pics/img/abcd.gif")
	assert.Equal(t, post.Likes, int64(123))
}

func TestCSVToPost(t *testing.T) {
	row := make([]string, 5)
	row[0] = "1234"
	row[1] = "title"
	row[2] = "url"
	row[3] = "abcd.gif"
	row[4] = "123"
	post := CSVToPost(row)
	assert.Equal(t, post.ID, int64(1234))
	assert.Equal(t, post.Title, "title")
	assert.Equal(t, post.URL, "url")
	assert.Equal(t, post.Image, "http://static.reaction.pics/img/abcd.gif")
	assert.Equal(t, post.Likes, int64(123))
}

func TestCorruptPost(t *testing.T) {
	row := make([]string, 5)
	post := CSVToPost(row)
	assert.Equal(t, post.ID, int64(0))
}

func TestInternalURL(t *testing.T) {
	post := Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 123}
	url := post.InternalURL()
	assert.Equal(t, url, "/post/1/title1")
}

func TestInternalURLLong(t *testing.T) {
	post := Post{1, strings.Repeat("a", 50), "url1", "http://static.reaction.pics/img/abcd.gif", 123}
	url := post.InternalURL()
	assert.Equal(t, url, "/post/1/"+strings.Repeat("a", 30))
}

func TestInitializeBoard(t *testing.T) {
	b := InitializeBoard()
	// When the lock is released, it should have data in it
	time.Sleep(time.Millisecond * 10)
	b.mut.RLock()
	b.mut.RUnlock()
	assert.True(t, len(b.Posts) > 0)
}

func TestAddPost(t *testing.T) {
	post := Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 123}
	board := NewBoard([]Post{})
	board.AddPost(post)
	assert.Equal(t, len(board.Posts), 1)
	assert.Equal(t, board.Posts[0].Title, "title1")

	board.AddPost(post)
	assert.Equal(t, len(board.Posts), 1)
}

func TestPostsToJSON(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 123}
	posts[1] = Post{2, "title2", "url2", "http://static.reaction.pics/img/abcd.gif", 124}
	board := NewBoard(posts)
	data := board.PostsToJSON()
	assert.Equal(t, len(*data), 2)

	post := (*data)[0]
	assert.Equal(t, post.ID, int64(1))
	assert.Equal(t, post.Title, "title1")
	assert.Equal(t, post.URL, "url1")
	assert.Equal(t, post.Image, "http://static.reaction.pics/img/abcd.gif")
	assert.Equal(t, post.Likes, int64(123))
	assert.Equal(t, post.InternalURL, "/post/1/title1")
}

func TestFilterBoard(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 123}
	posts[1] = Post{2, "title2", "url2", "http://static.reaction.pics/img/abcd.gif", 124}
	board := NewBoard(posts)
	newBoard := board.FilterBoard("title2")
	assert.Equal(t, len(newBoard.Posts), 1)
	assert.Equal(t, newBoard.Posts[0].ID, int64(2))
}

func TestLimitBoard(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 123}
	posts[1] = Post{2, "title2", "url2", "http://static.reaction.pics/img/abcd.gif", 124}
	board := NewBoard(posts)
	board.LimitBoard(1, 1)
	assert.Equal(t, len(board.Posts), 1)
	assert.Equal(t, board.Posts[0].ID, int64(2))

	board.LimitBoard(1, 3)
	assert.Equal(t, len(board.Posts), 0)

	board.LimitBoard(10, 1)
	assert.Equal(t, len(board.Posts), 0)
}

func TestSortPostsByLikes(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title3", "url3", "http://static.reaction.pics/img/abcd.gif", 123})
	board.AddPost(Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 121})
	board.AddPost(Post{2, "title2", "url2", "http://static.reaction.pics/img/abcd.gif", 122})
	board.SortPostsByLikes()
	assert.Equal(t, board.Posts[0].Likes, int64(123))
	assert.Equal(t, board.Posts[1].Likes, int64(122))
	assert.Equal(t, board.Posts[2].Likes, int64(121))
}

func TestRandomizePosts(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 121})
	board.AddPost(Post{2, "title2", "url2", "http://static.reaction.pics/img/abcd.gif", 122})
	board.AddPost(Post{3, "title3", "url3", "http://static.reaction.pics/img/abcd.gif", 123})
	randomized := false
	for i := 0; i < 10; i++ {
		// Technically a flaky test, but is expected to only fail in one out of 3^10 chances
		board.RandomizePosts()
		if board.Posts[0].ID != 3 {
			randomized = true
		}
	}
	if !randomized {
		assert.Fail(t, "Did not find randomized post")
	}
}

func TestReset(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title3", "url3", "http://static.reaction.pics/img/abcd.gif", 123})
	assert.Equal(t, len(board.Posts), 1)

	board.Reset()
	assert.Equal(t, len(board.Posts), 0)
}

func TestURLs(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title3", "url3", "http://static.reaction.pics/img/abcd.gif", 123})
	board.AddPost(Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 121})
	board.AddPost(Post{2, "title2", "url2", "http://static.reaction.pics/img/abcd.gif", 122})
	urls := board.URLs()
	assert.Equal(t, len(urls), 3)
	assert.Equal(t, urls[0], "/post/3/title3")
	assert.Equal(t, urls[1], "/post/1/title1")
	assert.Equal(t, urls[2], "/post/2/title2")
}

func TestKeywords(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title2", "url3", "http://static.reaction.pics/img/abcd.gif", 123})
	board.AddPost(Post{1, "title1", "url1", "http://static.reaction.pics/img/abcd.gif", 121})
	board.AddPost(Post{2, "title1 title2 title2", "url2", "http://static.reaction.pics/img/abcd.gif", 122})
	keywords := board.Keywords()
	assert.Equal(t, len(keywords), 2)
	assert.Equal(t, keywords[0], "title2")
	assert.Equal(t, keywords[1], "title1")
}

func TestKeywordsLong(t *testing.T) {
	board := NewBoard([]Post{})
	title := []string{}
	for x := 10000; x < 10100; x++ {
		title = append(title, strconv.FormatInt(int64(x), 10))
	}
	board.AddPost(Post{1, strings.Join(title, " "), "url1", "http://static.reaction.pics/img/abcd.gif", 121})
	keywords := board.Keywords()
	assert.Equal(t, len(keywords), MaxKeywords)
}
