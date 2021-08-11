package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/albertyw/reaction-pics/server"
	"github.com/joho/godotenv"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logWriter allows the logger the implement the io.Writer interface
type logWriter struct {
	l *zap.SugaredLogger
}

// Write receives bytes and writes to a logger
func (w logWriter) Write(p []byte) (int, error) {
	w.l.Info(string(p))
	return len(p), nil
}

func setupEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func setupRollbar() {
	rollbar.SetToken(os.Getenv("ROLLBAR_SERVER_TOKEN"))
	rollbar.SetEnvironment(os.Getenv("ENVIRONMENT"))
}

func getLogger() *zap.SugaredLogger {
	var config zap.Config
	if os.Getenv("ENVIRONMENT") == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		filename = "."
	}
	logFile := filepath.Join(filepath.Dir(filename), "logs", "app", "app.log")
	config.OutputPaths = []string{"stdout", logFile}
	logger, err := config.Build()
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
	server.Run(logger)
}
