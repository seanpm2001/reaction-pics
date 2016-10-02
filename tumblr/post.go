package tumblr

import (
	"github.com/MariaTerzieva/gotumblr"
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
