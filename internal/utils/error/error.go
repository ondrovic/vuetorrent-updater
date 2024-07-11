package error

import (
	"log"
	"updater/internal/utils/system"
)

// Log - Logs out the given error with a custom message. If no error is provided, it logs only the message as informational. 
func Log(err error, message string) {
	if err != nil {
		log.Fatalf("Error: %v\n%s", err, message)
	} else {
		log.Printf("%s", message) // Use log.Print for info messages if you don't want a newline at the end.
	}
	system.Exit(1)
}