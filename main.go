package main

import (
	"os"

	"github.com/albertyw/reaction-pics/server"
	"github.com/albertyw/reaction-pics/tumblr"
	_ "github.com/joho/godotenv/autoload"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stvp/rollbar"
)

func getNewRelicApp() newrelic.Application {
	newrelicKey := os.Getenv("NEWRELIC_KEY")
	config := newrelic.NewConfig("Reaction.pics", newrelicKey)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		panic(err)
	}
	return app
}

func setupRollbar() {
	rollbar.Token = os.Getenv("ROLLBAR_SERVER_TOKEN")
	rollbar.Environment = os.Getenv("ENVIRONMENT")
}

func main() {
	newrelicApp := getNewRelicApp()
	posts := make(chan tumblr.Post)
	go tumblr.GetPosts(posts)
	server.Run(posts, newrelicApp)
}
