package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus" // Logrus - Structured logger for Go

	"github.com/twibber/api/lib"    // Library containing project-specific configurations
	"github.com/twibber/api/router" // Router package for handling HTTP routes
)

// init function is called before the main function. Used for setting up logging, configuration, and Sentry.
func init() {
	// Setting up the log formatter with overridden environment colors, timestamp details, and padding.
	log.SetFormatter(&log.TextFormatter{
		EnvironmentOverrideColors: true,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           time.StampMilli, // Setting the timestamp format to include milliseconds
		PadLevelText:              true,            // Padding the text to align log levels
	})

	// In debug mode, potentially additional setup like database migration could occur.
	if lib.Config.Debug {
		// lib.MigrateDB() // Uncomment if database migration should occur at startup in debug mode
		log.SetLevel(log.DebugLevel) // Setting the log level to Debug if in debug mode
	}
}

// The main function starts the HTTP listener and logs fatal errors if the server fails to start.
func main() {
	// Starting the HTTP server and listening on the configured port.
	if err := router.Configure().Listen(fmt.Sprintf("%s:%s", "0.0.0.0", lib.Config.Port)); err != nil {
		log.WithError(err).WithField("port", lib.Config.Port).Fatal("failed to start listener")
	}
}
