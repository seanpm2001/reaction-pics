package main

import (
	_ "github.com/albertyw/devops-reactions-index/tumblr"
	"github.com/albertyw/devops-reactions-index/server"
)

func main() {
	// posts := tumblr.GetPosts()
	// tumblr.WritePostsToCSV(posts)
	server.Run()
}
