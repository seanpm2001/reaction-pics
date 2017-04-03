package tumblr

import (
	"github.com/MariaTerzieva/gotumblr"
	"testing"
)

func TestPost(t *testing.T) {
	post := Post{
		1234,
		"title",
		"url",
		"http://placehold.it/350x150",
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
}

func TestGoTumblrToPost(t *testing.T) {
	tumblrPost := gotumblr.TextPost{
		Title: "title",
	}
	tumblrPost.Id = 1234
	tumblrPost.Post_url = "url"
	tumblrPost.Body = "<img src=\"http://placehold.it/350x150\" />"
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
}

func TestCSVToPost(t *testing.T) {
	row := make([]string, 4)
	row[0] = "1234"
	row[1] = "title"
	row[2] = "url"
	row[3] = "http://placehold.it/350x150"
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
}

func TestCorruptPost(t *testing.T) {
	row := make([]string, 4)
	post := CSVToPost(row)
	if post.ID != 0 {
		t.Fail()
	}
}

func TestPostsToJSON(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1", "http://placehold.it/350x150"}
	posts[1] = Post{2, "title2", "url2", "http://placehold.it/350x150"}
	json := PostsToJSON(posts)
	expected := "[{\"id\":1,\"title\":\"title1\",\"url\":\"url1\",\"image\":\"http://placehold.it/350x150\"},{\"id\":2,\"title\":\"title2\",\"url\":\"url2\",\"image\":\"http://placehold.it/350x150\"}]"
	if json != expected {
		t.Fail()
	}
}

func TestSort(t *testing.T) {
	posts := make([]Post, 3)
	posts[0] = Post{3, "title3", "url3", "http://placehold.it/350x150"}
	posts[1] = Post{1, "title1", "url1", "http://placehold.it/350x150"}
	posts[2] = Post{2, "title2", "url2", "http://placehold.it/350x150"}
	posts = *SortPosts(&posts)
	if posts[0].ID != 3 {
		t.Fail()
	}
	if posts[1].ID != 2 {
		t.Fail()
	}
	if posts[2].ID != 1 {
		t.Fail()
	}
}
