package filedownloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	progress "updater/internal/utils/progressbar" // Assuming progress is a package, adjust import path accordingly.
)

var consoleMutex sync.Mutex // Mutex to synchronize access to the console output.

// DownloadFile downloads a file from the specified URL to the specified output directory.
func DownloadFile(name string, url string, outputDir string) (string, error) {
	// Create the output directory if it doesn't exist.
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Get the file name from the URL and construct the full path.
	fileName := filepath.Base(url)
	outputPath := filepath.Join(outputDir, fileName)

	// Create the file.
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close() // Ensure the file is closed when done.

	// Perform an HTTP GET request.
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: %v", resp.Status)
	}
	defer resp.Body.Close() // Ensure the response body is closed when done.

	// Initialize and start a progress bar for the download.
	barFileDownload := progress.NewDefaultBar(resp.ContentLength, fmt.Sprintf("Downloading: '%v'", name))
	defer barFileDownload.Finish() // Finish the progress bar when function exits.

	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	// Copy the content to the file and update the progress bar concurrently.
	_, err = io.Copy(io.MultiWriter(outFile, barFileDownload), resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Print a newline for cleaner console output when done downloading.
	fmt.Println()

	return outputPath, nil
}
