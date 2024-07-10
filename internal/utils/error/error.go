package error

import (
	"log"
	"updater/internal/utils/system"
)

// Log - Logs out the error
func Log(err error, message string) {
	log.Fatal(err, message)
	system.Exit(1)
}