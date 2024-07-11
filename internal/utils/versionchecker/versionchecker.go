package versionchecker

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
	"updater/internal/utils/http"
)

type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

type ResponseData struct {
	Assets  []Asset `json:"assets"`
	TagName string  `json:"tag_name"`
}

// IsNewVersion - Checks to see if there is a new version available.
func IsNewVersion(url, filepath string) (bool, string, string, string, error) {
	var responseData ResponseData
	err := http.GetAndParse(url, &responseData)
	if err != nil {
		return false, "", "", "", fmt.Errorf("failed to check for new version: %v", err)
	}

	// Use the first asset if available, otherwise default values
	assetName, assetUrl := getFirstAssetOrDefault("", "", responseData)
	releaseVersion := cleanTagName(responseData.TagName)
	installedVersion := getInstalledVersion(filepath)

	if releaseVersion != installedVersion {
		return true, assetName, assetUrl, releaseVersion, nil
	}
	return false, "", "", "", nil
}

// ReturnCleanedTagName - Removes a leading character if found from tag.
func cleanTagName(tag string) string {
	if len(tag) > 0 && unicode.IsLetter(rune(tag[0])) {
		return tag[1:]
	} else {
		return tag
	}
}

// getInstalledVersion retrieves the installed version from a file. If no valid version is found, it returns "0.0.0".
func getInstalledVersion(filepath string) string {
	if exists, _ := fileExists(filepath); !exists {
		return "0.0.0"
	}
	versionNum, err := readFile(filepath)
	if err != nil {
		return "0.0.0"
	}
	return versionNum
}

// Helper function to check if a file exists and return boolean result along with potential error.
func fileExists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// Helper function to read a file and return its content as string. Handles errors internally.
func readFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	fileContents := ""
	for i, line := range lines {
		if i == len(lines)-1 { // If this is the last element, don't add "\n"
			fileContents += line
		} else {
			fileContents += line + "\n"
		}
	}
	return fileContents, nil
}

// Helper function to get the first asset or default values if none are available.
func getFirstAssetOrDefault(name, url string, data ResponseData) (string, string) {
	if len(data.Assets) > 0 {
		return data.Assets[0].Name, data.Assets[0].DownloadUrl
	} else {
		return name, url
	}
}
