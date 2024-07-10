package system

import "os"

type Exitor interface {
	Exit(code int)
}

type OSExiter struct{}

func (e OSExiter) Exit(code int) {
	os.Exit(code)
}

var DefaultExitor Exitor = OSExiter{}

func Exit(code int) {
	DefaultExitor.Exit(code)
}