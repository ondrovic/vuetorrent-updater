package update

import (
	"fmt"
	"updater/internal/utils/filedownloader"
)

// GetUpdate - downloads the new update and returns it if successful, or error otherwise.
func GetUpdate(assetname string, url string, downloadpath string) (string, error) {
	file, err := filedownloader.DownloadFile(assetname, url, downloadpath)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	return file, nil
}