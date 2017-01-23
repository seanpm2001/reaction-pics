package main

import (
	"github.com/albertyw/devops-reactions-index/server"
	"github.com/albertyw/devops-reactions-index/tumblr"
	"os"
	"strings"
)

const readPostsFromTumblrEnv = "READ_POSTS_FROM_TUMBLR"

func getReadPostsFromTumblr() bool {
	var readPostsEnv = os.Getenv(readPostsFromTumblrEnv)
	if strings.ToLower(readPostsEnv) == "true" {
		return true
	}
	return false
}

func main() {
	readPosts := getReadPostsFromTumblr()
	posts := tumblr.GetPosts(readPosts)
	tumblr.WritePostsToCSV(posts)
	server.Run(posts)
}
