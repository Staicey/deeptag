package deepcrate

import (
	"fmt"
	"log"
	"os"
)

var deepcrateUrl string

type Completed struct {
	Items []struct {
		ID     string `json:"id"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
	} `json:"items"`
}

type WebHook struct {
	Event string `json:"event"`
	Data  struct {
		Artist       string `json:"artist"`
		Album        string `json:"album"`
		DownloadPath string `json:"download_path"`
	} `json:"data"`
}

func init() {
	if val, ok := os.LookupEnv("DEEPCRATE_URL"); ok == true {
		deepcrateUrl = fmt.Sprintf("%s/api/v1", val)
	} else {
		log.Fatal("[ERROR] Failed to load deepcrate environment variables: DEEPCRATE_URL")
		os.Exit(1)
	}

	log.Println("[INFO] Loaded deepcrate environment variables")
}
