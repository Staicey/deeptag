package deepcrate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/staicey/deeptag/util"
)

var (
	deepcrateUrl  string
	deepcrateUser string
	deepcratePass string
)

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

type auth struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
}

func init() {
	deepcrateUrl = util.LoadEnvVar("DEEPCRATE_URL") + "/api/v1"

	// Check authentication type
	var dcAuth auth
	code, out := util.Request("GET", fmt.Sprintf("%s/auth/info", deepcrateUrl), strings.NewReader(""), "application/json")
	if code != 200 {
		log.Fatal("[ERROR] Cannot connect to deepcrate, check DEEPCRATE_URL is configured correctly")
		os.Exit(1)
	}
	json.Unmarshal(out, &dcAuth)

	// Load the username/password
	if dcAuth.Enabled == true {
		log.Println("[INFO] Deepcrate authentication enabled, loading credentials")
		deepcrateUser = util.LoadEnvVar("DEEPCRATE_USER")
		deepcratePass = util.LoadEnvVar("DEEPCRATE_PASS")

		bearer := fmt.Sprintf("%s:%s", deepcrateUser, deepcratePass)
		elems := strings.Split(deepcrateUrl, "://")
		proto := elems[0]
		host := elems[1]
		deepcrateUrl = fmt.Sprintf("%s://%s@%s", proto, bearer, host)

		code, _ := util.Request("GET", fmt.Sprintf("%s/auth/me", deepcrateUrl), strings.NewReader(""), "application/json")
		if code != 200 {
			log.Fatal("[ERROR] Cannot authenticate with deepcrate, check DEEPCRATE_USER & DEEPCRATE_PASS are configured correctly")
			os.Exit(1)
		}
	}

	log.Println("[INFO] Loaded deepcrate environment variables")
}
