package util

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// Generic helper for making web requests
func Request(action string, url string, payload *strings.Reader, contentType string) (int, []byte) {
	req, _ := http.NewRequest(action, url, payload)

	req.Header.Add("Content-Type", contentType)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return res.StatusCode, body
}
