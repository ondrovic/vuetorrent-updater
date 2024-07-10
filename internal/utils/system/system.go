// system.go
package system

import "os"

// Exitor is an interface that defines the Exit function
type Exitor interface {
	Exit(code int)
}

// OSExiter is an implementation of Exitor that calls os.Exit
type OSExiter struct{}

// Exit exits the program with the given code
func (e OSExiter) Exit(code int) {
	os.Exit(code)
}

// DefaultExitor is the default implementation of Exitor
var DefaultExitor Exitor = OSExiter{}

// Exit calls the DefaultExitor's Exit method
func Exit(code int) {
	DefaultExitor.Exit(code)
}

// CheckError checks for an error and exits with the given code if an error is found
func CheckError(err error, code int) {
	if err != nil {
		Exit(code)
	}
}