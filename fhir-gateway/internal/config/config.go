package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var DBURL string

func Load() {
	DBURL = os.Getenv("DB_URL")
	if DBURL == "" {
		logrus.Fatal("DB_URL environment variable is required")
	}
}
