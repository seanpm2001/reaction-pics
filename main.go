package main

import (
	"github.com/albertyw/reaction-pics/server"
	"github.com/albertyw/reaction-pics/tumblr"
	_ "github.com/joho/godotenv/autoload"
	newrelic "github.com/newrelic/go-agent"
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

func getNewRelicApp() newrelic.Application {
	newrelicKey := os.Getenv("NEWRELIC_KEY")
	config := newrelic.NewConfig("Reaction.pics", newrelicKey)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		panic(err)
	}
	return app
}

func main() {
	newrelicApp := getNewRelicApp()
	readPosts := getReadPostsFromTumblr()
	posts := make(chan tumblr.Post)
	go tumblr.GetPosts(readPosts, posts)
	server.Run(posts, newrelicApp)
}
