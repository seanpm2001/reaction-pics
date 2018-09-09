package main

import (
	"os"

	"github.com/albertyw/reaction-pics/server"
	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/joho/godotenv"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rollbar/rollbar-go"
)

func setupEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
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

func setupRollbar() {
	rollbar.SetToken(os.Getenv("ROLLBAR_SERVER_TOKEN"))
	rollbar.SetEnvironment(os.Getenv("ENVIRONMENT"))
}

func main() {
	setupEnv()
	newrelicApp := getNewRelicApp()
	posts := make(chan tumblr.Post)
	go tumblr.GetPosts(posts)
	server.Run(posts, newrelicApp)
}
