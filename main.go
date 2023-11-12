package main

import (
	"fmt"
	"os"
	"time"

	sentryHook "github.com/chadsr/logrus-sentry" // Sentry hook for logrus
	"github.com/getsentry/sentry-go"             // Official Sentry Go SDK
	log "github.com/sirupsen/logrus"             // Logrus - Structured logger for Go

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

	// If debug mode is on, set the logger level to Debug.
	if lib.Config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// In debug mode, potentially additional setup like database migration could occur.
	if lib.Config.Debug {
		// lib.MigrateDB() // Uncomment if database migration should occur at startup in debug mode
		log.SetLevel(log.DebugLevel) // Setting the log level to Debug if in debug mode
	}

	// Configuring Sentry to only be used when not in debug mode.
	if !lib.Config.Debug {
		// Initializing Sentry with the project DSN, debugging options, and server name.
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              lib.Config.SentryDSN,
			Debug:            true,                  // Enabling Sentry debug mode
			AttachStacktrace: true,                  // Attaching stack trace to Sentry events
			SampleRate:       1,                     // Setting the sample rate for event reporting
			ServerName:       os.Getenv("HOSTNAME"), // Using the hostname from environment variable
		}); err != nil {
			log.Fatal(err) // Logging and exiting if Sentry initialization fails
		}

		// Adding the Sentry hook to logrus for capturing warnings, panics, fatals, and errors.
		log.AddHook(sentryHook.New([]log.Level{
			log.WarnLevel, log.PanicLevel, log.FatalLevel, log.ErrorLevel,
		}))
	}
}

// The main function starts the HTTP listener and logs fatal errors if the server fails to start.
func main() {
	// Starting the HTTP server and listening on the configured port.
	if err := router.Configure().Listen(fmt.Sprintf("%s:%s", "0.0.0.0", lib.Config.Port)); err != nil {
		log.WithError(err).WithField("port", lib.Config.Port).Fatal("failed to start listener")
	}
}
