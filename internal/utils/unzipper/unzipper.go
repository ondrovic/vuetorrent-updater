package unzipper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"updater/internal/utils/progressbar" // Assuming progressbar is in this package for simplicity. Adjust import path as necessary.
)

var (
	// Mutex to synchronize access to the console output
	consoleMutex sync.Mutex
)

// UnzipWithProgress extracts files from a ZIP archive with progress bar.
func UnzipWithProgress(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zipfile: %v", err)
	}
	defer r.Close()

	// Get total size of all files in the zip archive
	var totalSize int64 = 0
	for _, f := range r.File {
		totalSize += int64(f.UncompressedSize64)
	}

	barExtract := progressbar.NewDefaultBar(totalSize, fmt.Sprintf("Extracting '%s'", src))
	defer barExtract.Finish() // Defer the finish call to ensure it runs even if there's an error.

	// Synchronize output to prevent overlapping with other progress bars
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// Create directory
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Create the file's parent directories if necessary
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}

			// Create the file
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			defer outFile.Close() // Defer closing the file to ensure it runs even if there's an error copying contents.

			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in zip: %v", err)
			}
			defer rc.Close() // Defer closing the reader to ensure it runs even if there's an error copying contents.

			// Copy the file contents with progress
			_, err = io.Copy(io.MultiWriter(outFile, barExtract), rc)
			if err != nil {
				return fmt.Errorf("failed to copy file contents: %v", err)
			}
		}
	}

	fmt.Println()

	return nil
}