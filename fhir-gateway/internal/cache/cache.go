package cache

import (
	"fhir-gateway/internal/config"

	"github.com/sirupsen/logrus"
	"github.com/valkey-io/valkey-go"
)

func Connect() valkey.Client {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{config.CACHEURL},
	})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to cache server")
	}
	logrus.Info("Connected to cache server")
	return client
}
