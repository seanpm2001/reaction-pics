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

func duplicateChan(in <-chan tumblr.Post, out1, out2 chan<- tumblr.Post) {
	for p := range in {
		out1 <- p
		out2 <- p
	}
	close(out1)
	close(out2)
}

func main() {
	readPosts := getReadPostsFromTumblr()
	posts := make(chan tumblr.Post)
	posts1 := make(chan tumblr.Post)
	posts2 := make(chan tumblr.Post)
	go tumblr.GetPosts(readPosts, posts)
	go duplicateChan(posts, posts1, posts2)
	go tumblr.WritePostsToCSV(posts1)
	server.Run(posts2)
}
