package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	Logger.SetLevel(logrus.InfoLevel)
}
