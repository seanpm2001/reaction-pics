package tumblr

import (
	"encoding/json"
	"github.com/MariaTerzieva/gotumblr"
	"sort"
	"strconv"
)

// Post is a representation of a single tumblr post
type Post struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// GoTumblrToPost converts a gotumblr.TextPost into a Post
func GoTumblrToPost(tumblrPost *gotumblr.TextPost) *Post {
	post := Post{
		ID:    tumblrPost.Id,
		Title: tumblrPost.Title,
		URL:   tumblrPost.Post_url,
	}
	return &post
}

// CSVToPost converts a CSV row into a Post
func CSVToPost(row []string) *Post {
	id, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		id = 0
	}
	post := Post{
		ID:    id,
		Title: row[1],
		URL:   row[2],
	}
	return &post
}

// PostsToJSON converts a Post into a JSON string
func PostsToJSON(posts []Post) string {
	jsonPosts, _ := json.Marshal(posts)
	return string(jsonPosts)
}

// SortPosts sorts Posts in reverse ID order
func SortPosts(posts *[]Post) *[]Post {
	sort.Sort(SortByID(*posts))
	sort.Sort(sort.Reverse(SortByID(*posts)))
	return posts
}

// SortByID is an interface for Sorting
type SortByID []Post

func (a SortByID) Len() int           { return len(a) }
func (a SortByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
