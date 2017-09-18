package tumblr

import (
	"strconv"
	"strings"
	"testing"

	"github.com/MariaTerzieva/gotumblr"
)

func TestPost(t *testing.T) {
	post := Post{
		1234,
		"title",
		"url",
		"http://placehold.it/350x150",
		123,
	}
	if post.ID != 1234 {
		t.Fail()
	}
	if post.Title != "title" {
		t.Fail()
	}
	if post.URL != "url" {
		t.Fail()
	}
	if post.Image != "http://placehold.it/350x150" {
		t.Fail()
	}
	if post.Likes != 123 {
		t.Fail()
	}
}

func TestGoTumblrToPost(t *testing.T) {
	tumblrPost := gotumblr.TextPost{
		Title: "title",
	}
	tumblrPost.Id = 1234
	tumblrPost.Post_url = "url"
	tumblrPost.Body = "<img src=\"http://placehold.it/350x150\" />"
	tumblrPost.Note_count = 123
	post := GoTumblrToPost(&tumblrPost)
	if post.ID != 1234 {
		t.Fail()
	}
	if post.Title != "title" {
		t.Fail()
	}
	if post.URL != "url" {
		t.Fail()
	}
	if post.Image != "http://placehold.it/350x150" {
		t.Fail()
	}
	if post.Likes != 123 {
		t.Fail()
	}
}

func TestCSVToPost(t *testing.T) {
	row := make([]string, 5)
	row[0] = "1234"
	row[1] = "title"
	row[2] = "url"
	row[3] = "http://placehold.it/350x150"
	row[4] = "123"
	post := CSVToPost(row)
	if post.ID != 1234 {
		t.Fail()
	}
	if post.Title != "title" {
		t.Fail()
	}
	if post.URL != "url" {
		t.Fail()
	}
	if post.Image != "http://placehold.it/350x150" {
		t.Fail()
	}
	if post.Likes != 123 {
		t.Fail()
	}
}

func TestCorruptPost(t *testing.T) {
	row := make([]string, 5)
	post := CSVToPost(row)
	if post.ID != 0 {
		t.Fail()
	}
}

func TestInternalURL(t *testing.T) {
	post := Post{1, "title1", "url1", "http://placehold.it/350x150", 123}
	url := post.InternalURL()
	if url != "/post/1/title1" {
		t.Fail()
	}
}

func TestInternalURLLong(t *testing.T) {
	post := Post{1, strings.Repeat("a", 50), "url1", "http://placehold.it/350x150", 123}
	url := post.InternalURL()
	if url != "/post/1/"+strings.Repeat("a", 30) {
		t.Fail()
	}
}

func TestPostsToJSON(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1", "http://placehold.it/350x150", 123}
	posts[1] = Post{2, "title2", "url2", "http://placehold.it/350x150", 124}
	board := NewBoard(posts)
	data := board.PostsToJSON()
	if len(*data) != 2 {
		t.Fail()
	}
	post := (*data)[0]
	if post.ID != 1 {
		t.Fail()
	}
	if post.Title != "title1" {
		t.Fail()
	}
	if post.URL != "url1" {
		t.Fail()
	}
	if post.Image != "http://placehold.it/350x150" {
		t.Fail()
	}
	if post.Likes != 123 {
		t.Fail()
	}
	if post.InternalURL != "/post/1/title1" {
		t.Fail()
	}
}

func TestFilterBoard(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1", "http://placehold.it/350x150", 123}
	posts[1] = Post{2, "title2", "url2", "http://placehold.it/350x150", 124}
	board := NewBoard(posts)
	newBoard := board.FilterBoard("title2")
	if len(newBoard.Posts) != 1 {
		t.Fail()
	}
	if newBoard.Posts[0].ID != 2 {
		t.Fail()
	}
}

func TestLimitBoard(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1", "http://placehold.it/350x150", 123}
	posts[1] = Post{2, "title2", "url2", "http://placehold.it/350x150", 124}
	board := NewBoard(posts)
	board.LimitBoard(1, 1)
	if len(board.Posts) != 1 {
		t.Fail()
	}
	if board.Posts[0].ID != 2 {
		t.Fail()
	}
	board.LimitBoard(1, 3)
	if len(board.Posts) != 0 {
		t.Fail()
	}
	board.LimitBoard(10, 1)
	if len(board.Posts) != 0 {
		t.Fail()
	}
}

func TestGetPostByID(t *testing.T) {
	posts := make([]Post, 1)
	posts[0] = Post{1, "title1", "url1", "http://placehold.it/350x150", 123}
	board := NewBoard(posts)
	post := board.GetPostByID(1)
	if post == nil {
		t.Fail()
	}
	if post.Title != "title1" {
		t.Fail()
	}
	post = board.GetPostByID(2)
	if post != nil {
		t.Fail()
	}
}

func TestSortPostsByID(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title3", "url3", "http://placehold.it/350x150", 123})
	board.AddPost(Post{1, "title1", "url1", "http://placehold.it/350x150", 124})
	board.AddPost(Post{2, "title2", "url2", "http://placehold.it/350x150", 125})
	board.SortPostsByID()
	if board.Posts[0].ID != 3 {
		t.Fail()
	}
	if board.Posts[1].ID != 2 {
		t.Fail()
	}
	if board.Posts[2].ID != 1 {
		t.Fail()
	}
}

func TestSortPostsByLikes(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title3", "url3", "http://placehold.it/350x150", 123})
	board.AddPost(Post{1, "title1", "url1", "http://placehold.it/350x150", 121})
	board.AddPost(Post{2, "title2", "url2", "http://placehold.it/350x150", 122})
	board.SortPostsByLikes()
	if board.Posts[0].Likes != 123 {
		t.Fail()
	}
	if board.Posts[1].Likes != 122 {
		t.Fail()
	}
	if board.Posts[2].Likes != 121 {
		t.Fail()
	}
}

func TestReset(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title3", "url3", "http://placehold.it/350x150", 123})
	if len(board.Posts) != 1 {
		t.Fail()
	}
	board.Reset()
	if len(board.Posts) != 0 {
		t.Fail()
	}
}

func TestURLs(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title3", "url3", "http://placehold.it/350x150", 123})
	board.AddPost(Post{1, "title1", "url1", "http://placehold.it/350x150", 121})
	board.AddPost(Post{2, "title2", "url2", "http://placehold.it/350x150", 122})
	urls := board.URLs()
	if len(urls) != 3 {
		t.Fail()
	}
	if urls[0] != "/post/3/title3" {
		t.Fail()
	}
	if urls[1] != "/post/1/title1" {
		t.Fail()
	}
	if urls[2] != "/post/2/title2" {
		t.Fail()
	}
}

func TestKeywords(t *testing.T) {
	board := NewBoard([]Post{})
	board.AddPost(Post{3, "title2", "url3", "http://placehold.it/350x150", 123})
	board.AddPost(Post{1, "title1", "url1", "http://placehold.it/350x150", 121})
	board.AddPost(Post{2, "title1 title2 title2", "url2", "http://placehold.it/350x150", 122})
	keywords := board.Keywords()
	if len(keywords) != 2 {
		t.Fail()
	}
	if keywords[0] != "title2" {
		t.Fail()
	}
	if keywords[1] != "title1" {
		t.Fail()
	}
}

func TestKeywordsLong(t *testing.T) {
	board := NewBoard([]Post{})
	title := []string{}
	for x := 10000; x < 10100; x++ {
		title = append(title, strconv.FormatInt(int64(x), 10))
	}
	board.AddPost(Post{1, strings.Join(title, " "), "url1", "http://placehold.it/350x150", 121})
	keywords := board.Keywords()
	if len(keywords) != MaxKeywords {
		t.Fail()
	}
}

func TestGetImageFromPostBody(t *testing.T) {
	body := "<img src=\"img.gif\">"
	image := getImageFromPostBody(body)
	if image != "img.gif" {
		t.Fail()
	}
}

func TestGetImageFromPostBodyNotFound(t *testing.T) {
	body := "<div></div>"
	image := getImageFromPostBody(body)
	if image != "" {
		t.Fail()
	}
}
