package main

import (
	"github.com/albertyw/devops-reactions-index/server"
	"github.com/albertyw/devops-reactions-index/tumblr"
)

func main() {
	posts := tumblr.GetPosts()
	tumblr.WritePostsToCSV(posts)
	server.Run()
}
