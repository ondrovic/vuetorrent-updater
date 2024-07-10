package filedownloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"updater/internal/utils/progressbar"
)

var (
    // Mutex to synchronize access to the console output
    consoleMutex sync.Mutex
)

// DownloadFile downloads a file from the specified URL to the specified output directory.
func DownloadFile(name string, url string, outputDir string) (string, error) {
    // Create the output directory if it doesn't exist
    err := os.MkdirAll(outputDir, os.ModePerm)
    if err != nil {
        return "", fmt.Errorf("failed to create output directory: %v", err)
    }

    // Get the file name from the URL
    fileName := filepath.Base(url)
    outputPath := filepath.Join(outputDir, fileName)

    // Create the file
    outFile, err := os.Create(outputPath)
    if err != nil {
        return "", fmt.Errorf("failed to create file: %v", err)
    }
    defer outFile.Close()

    // Create a new HTTP request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", fmt.Errorf("failed to create HTTP request: %v", err)
    }

    // Perform the request
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to perform HTTP request: %v", err)
    }
    defer resp.Body.Close()

    // Check if the request was successful
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to download file: %v", resp.Status)
    }


    barFileDownload := progressbar.NewDefaultBar(
        resp.ContentLength,
        fmt.Sprintf("Downloading: '%v'", name),
    )

    // Synchronize output to prevent overlapping with other progress bars
    consoleMutex.Lock()
    defer consoleMutex.Unlock()

    // Copy the content to the file and update the progress bar
    _, err = io.Copy(io.MultiWriter(outFile, barFileDownload), resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to save file: %v", err)
    }

    // Finish the progress bar
    barFileDownload.Finish()

    fmt.Println()

    return outputPath, nil
}