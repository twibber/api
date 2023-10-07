package main

import (
	"fmt"
	"os"
	"time"

	sentryHook "github.com/chadsr/logrus-sentry"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"

	"github.com/twibber/api/lib"
	"github.com/twibber/api/router"
)

func init() {
	// lookin fancy
	log.SetFormatter(&log.TextFormatter{
		EnvironmentOverrideColors: true,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           time.StampMilli,
		PadLevelText:              true,
	})

	if lib.Config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// set logger level
	if lib.Config.Debug {
		// lib.MigrateDB()
		log.SetLevel(log.DebugLevel)
	}

	// use sentry if in prod
	if !lib.Config.Debug {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              lib.Config.SentryDSN,
			Debug:            true,
			AttachStacktrace: true,
			SampleRate:       1,
			ServerName:       os.Getenv("HOSTNAME"), // hostname is the pod name in kubernetes 😉
		}); err != nil {
			log.Fatal(err)
		}

		log.AddHook(sentryHook.New([]log.Level{log.WarnLevel, log.PanicLevel, log.FatalLevel, log.ErrorLevel}))
	}
}

func main() {
	// start http listener
	if err := router.Configure().Listen(fmt.Sprintf("%s:%s", "0.0.0.0", lib.Config.Port)); err != nil {
		log.WithError(err).WithField("port", lib.Config.Port).Fatal("failed to start listener")
	}
}
