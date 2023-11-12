package lib

import (
	"fmt"
	log "github.com/sirupsen/logrus"        // Logrus for structured logging
	"github.com/twibber/api/lib/gormLogger" // Custom GORM logger from the lib package
	"github.com/twibber/api/models"         // Data models for the application
	"gorm.io/driver/postgres"               // GORM driver for PostgreSQL
	"gorm.io/gorm"                          // GORM ORM library
	"reflect"                               // Standard library package for runtime reflection
	"strings"                               // Standard library package for string manipulation
)

var (
	DB *gorm.DB // DB is a global variable for the database connection
)

// init function establishes a connection to the database and logs the event.
func init() {
	// Opens a new database connection using the provided credentials and configuration.
	if conn, err := gorm.Open(postgres.Open(fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		Config.DBUsername, Config.DBPassword, Config.DBHost, Config.DBPort, Config.DBName)), &gorm.Config{
		Logger:               gormLogger.New(), // Uses the custom GORM logger
		FullSaveAssociations: true,             // Enables automatic saving of associated entities
	}); err != nil {
		// If connection fails, log the error and stop the application
		log.WithError(err).Fatal("could not connect to database")
	} else {
		// Log the success of the database connection attempt
		log.WithFields(log.Fields{
			"user":     Config.DBUsername,
			"password": strings.Repeat("*", len(Config.DBPassword)), // Masks the password for security
			"host":     Config.DBHost + ":" + Config.DBPort,
			"database": Config.DBName,
		}).Info("initiated database connection")
		DB = conn // Set the global DB variable to the connection
	}
}

// MigrateDB applies the auto migrations for the database models.
func MigrateDB() {
	// AutoMigrate will create or update database tables according to the models
	if err := DB.Migrator().AutoMigrate(models.Models...); err != nil {
		// Logs and stops the application if migration fails
		log.WithError(err).Fatal("could not migrate database")
	}

	// Retrieves and logs the names of all migrated models
	modelNames := make([]string, 0)
	for _, n := range models.Models {
		modelNames = append(modelNames, reflect.TypeOf(n).Elem().Name())
	}

	// Logs the successful migration of all models
	log.WithField("models", modelNames).Info("migrated all database models")
}
