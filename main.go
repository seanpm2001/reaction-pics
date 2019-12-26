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

func getNewRelicApp(logger *zap.SugaredLogger) newrelic.Application {
	newrelicKey := os.Getenv("NEWRELIC_KEY")
	config := newrelic.NewConfig("Reaction.pics", newrelicKey)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		logger.Fatal(err)
	}
	return app
}

func setupRollbar() {
	rollbar.SetToken(os.Getenv("ROLLBAR_SERVER_TOKEN"))
	rollbar.SetEnvironment(os.Getenv("ENVIRONMENT"))
}

func getLogger() *zap.SugaredLogger {
	var logger *zap.Logger
	var err error
	if os.Getenv("ENVIRONMENT") == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	sugaredLogger := logger.Sugar()
	return sugaredLogger
}

func main() {
	setupEnv()
	setupRollbar()
	logger := getLogger()
	defer logger.Sync()
	newrelicApp := getNewRelicApp(logger)
	server.Run(newrelicApp, logger)
}
