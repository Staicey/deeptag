package wrtag

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/staicey/deeptag/web"
)

var (
	wrUrl     string
	wrToken   string
	wrOp      string
	wrPattern string = "%s://%s@%s/op/%s"
)

func Import(filepath string) (int, []byte) {
	payload := strings.NewReader(fmt.Sprintf("path=/download/%s", filepath))
	code, out := web.Request("POST", wrUrl, payload, "application/x-www-form-urlencoded")

	log.Println("[DEBUG] import:", code, string(out))
	log.Println("[DEBUG] import: payload:", payload)

	log.Println("[INFO] Forward request to:", wrUrl)

	return code, out
}

// Load environment vars & setup endpoints
func init() {
	// Load the token from env
	if val, ok := os.LookupEnv("WRTAG_TOKEN"); ok == true {
		wrToken = ":" + val
	} else {
		log.Fatal("[ERROR] Failed to load wrtag environment variables: WRTAG_TOKEN")
		os.Exit(1)
	}

	// Load the operation from env
	if val, ok := os.LookupEnv("WRTAG_OP"); ok == true {
		wrOp = strings.ToLower(val)
		if !slices.Contains([]string{"copy", "move", "reflink"}, wrOp) {
			log.Fatal("[ERROR] Failed to validate WRTAG_OP, should be any of the following: [ copy | move | reflink ]")
			os.Exit(1)
		}
	} else {
		log.Fatal("[ERROR] Failed to load wrtag environment variables: WRTAG_OP")
		os.Exit(1)
	}

	if val, ok := os.LookupEnv("WRTAG_URL"); ok == true {
		elems := strings.Split(val, "://")
		proto := elems[0]
		hostname := elems[1]

		wrUrl = fmt.Sprintf(wrPattern, proto, wrToken, hostname, wrOp)
	} else {
		log.Fatal("[ERROR] Failed to load wrtag environment variables: RTAG_URL")
		os.Exit(1)
	}

	log.Println("[INFO] Loaded wrtag environment variables")
}
