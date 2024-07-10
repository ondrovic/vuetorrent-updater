package update

import (
    "fmt"

	"updater/internal/utils/http"
	"updater/internal/utils/filedownloader"
)

type Asset struct {
	Name string `json:"name"`
	Url string `json:"browser_download_url"`
}

type ResponseData struct {
	Assets []Asset `json:"assets"`
	TagName string `json:"tag_name"`
}

func RunCommand(outputDir string) (string, error) {
	url := "http://api.github.com/repos/VueTorrent/VueTorrent/releases/latest"

	// Create an instance of your custom data structure
	var data ResponseData

	// Call the helper function and pass the instance
	err := http.GetAndParse(url, &data)
	if err != nil {
		return "", fmt.Errorf("failed to get and parse response: %v", err)
	}
	tagName := data.TagName
	downloadUrl := ""
	assetName := ""
	if len(data.Assets) > 0 {
		assetName = data.Assets[0].Name
		downloadUrl = data.Assets[0].Url
	}

	fmt.Printf("Version: %s\n", tagName)

	// Download the file
	downloadedFilePath, err := filedownloader.DownloadFile(assetName, downloadUrl, outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}

	return downloadedFilePath, nil
}