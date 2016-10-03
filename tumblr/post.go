package tumblr

import (
	"github.com/MariaTerzieva/gotumblr"
	"sort"
	"strconv"
)

type Post struct {
	Id    int64
	Title string
	Url   string
}

func TumblrToPost(tumblrPost *gotumblr.TextPost) *Post {
	post := Post{
		Id:    tumblrPost.Id,
		Title: tumblrPost.Title,
		Url:   tumblrPost.Post_url,
	}
	return &post
}

func CSVToPost(row []string) *Post {
	id, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		id = 0
	}
	post := Post{
		Id:    id,
		Title: row[1],
		Url:   row[2],
	}
	return &post
}

func SortPosts(posts *[]Post) *[]Post {
	sort.Sort(SortById(*posts))
	sort.Sort(sort.Reverse(SortById(*posts)))
	return posts
}

type SortById []Post

func (a SortById) Len() int           { return len(a) }
func (a SortById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortById) Less(i, j int) bool { return a[i].Id < a[j].Id }
