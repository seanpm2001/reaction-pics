package tumblr

import (
	"github.com/MariaTerzieva/gotumblr"
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
