package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var DBURL string
var CACHEURL string

func Load() {
	DBURL = os.Getenv("DB_URL")
	if DBURL == "" {
		logrus.Fatal("DB_URL environment variable is required")
	}
	CACHEURL = os.Getenv("CACHE_URL")
	if CACHEURL == "" {
		logrus.Fatal("CACHE_URL environment variable is required")
	}
}
