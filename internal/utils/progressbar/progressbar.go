package progressbar

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

// NewDefaultBar creates a new default progress bar with specified size and description
func NewDefaultBar(size int64, description string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("[green]%s[reset]", description)),
		progressbar.OptionFullWidth(),
	)

	return bar
}