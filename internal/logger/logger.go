package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
