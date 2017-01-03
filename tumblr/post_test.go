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
}

func TestGoTumblrToPost(t *testing.T) {
	tumblrPost := gotumblr.TextPost{
		Title: "title",
	}
	tumblrPost.Id = 1234
	tumblrPost.Post_url = "url"
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
}

func TestCSVToPost(t *testing.T) {
	row := make([]string, 3)
	row[0] = "1234"
	row[1] = "title"
	row[2] = "url"
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
}

func TestCorruptPost(t *testing.T) {
	row := make([]string, 3)
	post := CSVToPost(row)
	if post.ID != 0 {
		t.Fail()
	}
}

func TestPostsToJSON(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1"}
	posts[1] = Post{2, "title2", "url2"}
	json := PostsToJSON(posts)
	expected := "[{\"id\":1,\"title\":\"title1\",\"url\":\"url1\"},{\"id\":2,\"title\":\"title2\",\"url\":\"url2\"}]"
	if json != expected {
		t.Fail()
	}
}

func TestSort(t *testing.T) {
	posts := make([]Post, 3)
	posts[0] = Post{3, "title3", "url3"}
	posts[1] = Post{1, "title1", "url1"}
	posts[2] = Post{2, "title2", "url2"}
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
