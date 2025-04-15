package logService

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogrus() {
	if os.Getenv("APP_ENV") != "production" {
		logrus.SetLevel(logrus.TraceLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
}
