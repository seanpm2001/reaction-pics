package main

import (
	"github.com/albertyw/devops-reactions-index/tumblr"
)

func main() {
	posts := tumblr.GetPosts()
	tumblr.WritePostsToCSV(posts)
}
