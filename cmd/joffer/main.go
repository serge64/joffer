package main

import (
	"guthub.com/serge64/joffer/internal/config"
	"guthub.com/serge64/joffer/internal/server"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	formatter := &logrus.TextFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true

	logrus.SetFormatter(formatter)

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("No .env file found")
	}
}

func main() {
	config := config.New()

	if err := server.Start(config); err != nil {
		logrus.Error(err)
	}
}
