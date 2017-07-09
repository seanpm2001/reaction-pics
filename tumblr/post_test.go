package tumblr

import (
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	"testing"
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

func TestPostsToJSON(t *testing.T) {
	posts := make([]Post, 2)
	posts[0] = Post{1, "title1", "url1", "http://placehold.it/350x150", 123}
	posts[1] = Post{2, "title2", "url2", "http://placehold.it/350x150", 124}
	json := PostsToJSON(posts)
	expected := "[{\"id\":1,\"title\":\"title1\",\"url\":\"url1\",\"image\":\"http://placehold.it/350x150\",\"likes\":123,\"internalURL\":\"/post/1/title1\"},{\"id\":2,\"title\":\"title2\",\"url\":\"url2\",\"image\":\"http://placehold.it/350x150\",\"likes\":124,\"internalURL\":\"/post/2/title2\"}]"
	fmt.Println(json)
	fmt.Println(expected)
	if json != expected {
		t.Fail()
	}
}

func TestSortPostsByID(t *testing.T) {
	posts := make([]Post, 3)
	posts[0] = Post{3, "title3", "url3", "http://placehold.it/350x150", 123}
	posts[1] = Post{1, "title1", "url1", "http://placehold.it/350x150", 124}
	posts[2] = Post{2, "title2", "url2", "http://placehold.it/350x150", 125}
	posts = *SortPostsByID(&posts)
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

func TestSortPostsByLikes(t *testing.T) {
	posts := make([]Post, 3)
	posts[0] = Post{3, "title3", "url3", "http://placehold.it/350x150", 123}
	posts[1] = Post{1, "title1", "url1", "http://placehold.it/350x150", 121}
	posts[2] = Post{2, "title2", "url2", "http://placehold.it/350x150", 122}
	posts = *SortPostsByLikes(&posts)
	if posts[0].Likes != 123 {
		t.Fail()
	}
	if posts[1].Likes != 122 {
		t.Fail()
	}
	if posts[2].Likes != 121 {
		t.Fail()
	}
}
