package lib

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twibber/api/lib/gormLogger"
	"github.com/twibber/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

var (
	DB *gorm.DB
)

func init() {
	if conn, err := gorm.Open(postgres.Open(fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", Config.DBUsername, Config.DBPassword, Config.DBHost, Config.DBPort, Config.DBName)), &gorm.Config{
		Logger:               gormLogger.New(),
		FullSaveAssociations: true,
	}); err != nil {
		log.WithError(err).Fatal("could not connect to database")
	} else {
		log.WithFields(log.Fields{
			"user":     Config.DBUsername,
			"password": strings.Repeat("*", len(Config.DBPassword)), "host": Config.DBHost + ":" + Config.DBPort, "database": Config.DBName,
		}).Info("initiated database connection")
		DB = conn
	}
}

func MigrateDB() {
	if err := DB.Migrator().AutoMigrate(models.Models...); err != nil {
		log.WithError(err).Fatal("could not migrate database")
	}

	// just for the shits and giggles
	modelNames := make([]string, 0)
	for _, n := range models.Models {
		modelNames = append(modelNames, reflect.TypeOf(n).Elem().Name())
	}

	log.WithField("models", modelNames).Info("migrated all database models")
}
