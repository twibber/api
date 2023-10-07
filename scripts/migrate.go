package main

import (
	"github.com/twibber/api/lib"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("migrating database")
	lib.MigrateDB()
}
