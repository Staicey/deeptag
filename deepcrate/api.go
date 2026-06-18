package deepcrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/staicey/deeptag/web"
)

// Add an album to wishlist
func AddToWishlist(event *WebHook) {
	payload := strings.NewReader(fmt.Sprintf(`{
    "artist": "%s",
    "title": "%s",
    "type": "album",
    "year": 1,
    "mbid": ""
}`, event.Data.Artist, event.Data.Album))

	url := fmt.Sprintf("%s/wishlist", deepcrateUrl)
	code, out := web.Request("POST", url, payload, "application/json")

	log.Println("[DEBUG] create:", code, string(out))

	log.Println("[INFO] Adding to wishlist:", event.Data.Artist, "-", event.Data.Album)
}

// Get the download ID of a recently completed album
func GetDownloadID(target *WebHook) string {
	payload := strings.NewReader("")
	url := fmt.Sprintf("%s/downloads/completed?limit=50&offset=0", deepcrateUrl)
	_, resp := web.Request("GET", url, payload, "application/json")

	var downloads Completed
	json.NewDecoder(bytes.NewReader(resp)).Decode(&downloads)

	for _, item := range downloads.Items {
		if item.Album == target.Data.Album && item.Artist == target.Data.Artist {
			return item.ID
		}
	}

	return ""
}

// Remove a download ID from the database
func DeleteDownloadRecord(id string) {
	log.Println("[INFO] Deleting:", id, "from deepcrate")
	payload := strings.NewReader(fmt.Sprintf(`{
  "ids": [
    "%s"
  ]
}`, id))

	url := fmt.Sprintf("%s/downloads/", deepcrateUrl)
	code, out := web.Request("DELETE", url, payload, "application/json")

	log.Println("[DEBUG] delete:", code, string(out))
}

// Trigger an internal job
func Trigger(job string) {
	payload := strings.NewReader("")
	url := fmt.Sprintf("%s/jobs/%s/trigger", deepcrateUrl, job)
	code, out := web.Request("POST", url, payload, "application/json")

	log.Println("[DEBUG] trigger:", code, string(out))
}
