package main

import (
	"encoding/json"
	"github.com/MariaTerzieva/gotumblr"
	"testing"
)

func makeMockPost() gotumblr.PostsResponse {
	mockBlogInfo := gotumblr.BlogInfo{
		Title: "title",
		Posts: 10,
		Name:  "name",
		Url:   "url",
	}
	mockTextPost := gotumblr.TextPost{
		Title: "title",
		Body:  "Body",
	}
	mockTextPostJson, _ := json.Marshal(&mockTextPost)
	mockTextPostJsonRawMessage := *(*json.RawMessage)(&mockTextPostJson)
	mockPostsResponse := gotumblr.PostsResponse{
		Blog:        mockBlogInfo,
		Posts:       []json.RawMessage{mockTextPostJsonRawMessage},
		Total_posts: 5,
	}
	return mockPostsResponse
}

func TestParsePosts(t *testing.T) {

}
