package filedeleter

import (
	"fmt"
	"os"
	"sync"

	progress "updater/internal/utils/progressbar"
)

var (
	// Mutex to synchronize access to the console output
	consoleMutex sync.Mutex
)

// DeleteFileWithProgress deletes a file and shows progress using a progress bar.
func DeleteFileWithProgress(filepath string) error {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("failed to get file information: %v", err)
	}

	barDelete := progress.NewDefaultBar(
		fileInfo.Size(),
		fmt.Sprintf("Deleting '%s'", filepath),
	)

	// Synchronize output to prevent overlapping with other progress bars
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	// Delete the file
	err = os.Remove(filepath)
	if err != nil {
		barDelete.Finish() // Finish progress bar in case of error
		return fmt.Errorf("failed to delete file: %v", err)
	}

	// Finish deletion progress bar
	barDelete.Finish()

	fmt.Println()

	return nil
}

// DeleteFileWithProgress deletes a file without using a progress bar.
func DeleteFile(filepath string) error {
	// Synchronize output to prevent overlapping with other operations
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	// Delete the file
	err := os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	fmt.Printf("\nDeleted file '%s'\n", filepath)
	return nil
}