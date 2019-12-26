package main

import (
	"os"

	"github.com/albertyw/reaction-pics/server"
	"github.com/joho/godotenv"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
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

func getLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugaredLogger := logger.Sugar()
	return sugaredLogger
}

func main() {
	setupEnv()
	setupRollbar()
	logger := getLogger()
	newrelicApp := getNewRelicApp()
	server.Run(newrelicApp, logger)
}
