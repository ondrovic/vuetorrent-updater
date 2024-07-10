package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// GetAndParse performs an HTTP GET request to the specified URL and parses the JSON response into the provided data structure.
func GetAndParse(url string, target interface{}) error {
    // Perform the HTTP GET request
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Check if the request was successful
    if resp.StatusCode != http.StatusOK {
        return errors.New(resp.Status)
    }

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // Parse the JSON response
    err = json.Unmarshal(body, target)
    if err != nil {
        return err
    }

    return nil
}