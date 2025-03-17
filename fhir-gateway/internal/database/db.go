package database

import (
	"fhir-gateway/internal/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.DBURL), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	// err = db.AutoMigrate(&Patient{}, &APIKey{})
	// if err != nil {
	// 	logrus.WithError(err).Fatal("Failed to run migrations")
	// }
	logrus.Info("Connected to database")
	return db
}

type Patient struct {
	ID   string `gorm:"primaryKey"`
	Data []byte `gorm:"type:jsonb"`
}

type APIKey struct {
	Key    string `gorm:"primaryKey"`
	Active bool
}
