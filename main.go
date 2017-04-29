package main

import (
	"github.com/albertyw/reaction-pics/server"
	"github.com/albertyw/reaction-pics/tumblr"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
)

const readPostsFromTumblrEnv = "READ_POSTS_FROM_TUMBLR"

func getReadPostsFromTumblr() bool {
	readPostsEnv := os.Getenv(readPostsFromTumblrEnv)
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
