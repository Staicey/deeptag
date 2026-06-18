package wrtag

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/staicey/deeptag/util"
)

var (
	wrUrl     string
	wrToken   string
	wrOp      string
	wrDlDir   string
	wrPattern string = "%s://%s@%s/op/%s"
)

func Import(filepath string) (int, []byte) {
	payload := strings.NewReader(fmt.Sprintf("path=%s/%s", wrDlDir, filepath))
	code, out := util.Request("POST", wrUrl, payload, "application/x-www-form-urlencoded")

	log.Println("[DEBUG] import:", code, string(out))
	log.Println("[DEBUG] import: payload:", payload)

	log.Println("[INFO] Forward request to:", wrUrl)

	return code, out
}

// Load environment vars & setup endpoints
func init() {
	wrToken = ":" + util.LoadEnvVar("WRTAG_TOKEN")
	wrDlDir = util.LoadEnvVar("WRTAG_DL_DIR")

	wrOp = util.LoadEnvVar("WRTAG_OP")
	// Validate the op is within a set
	if !slices.Contains([]string{"copy", "move", "reflink"}, wrOp) {
		log.Fatal("[ERROR] Failed to validate WRTAG_OP, should be any of the following: [ copy | move | reflink ]")
		os.Exit(1)
	}

	elems := strings.Split(util.LoadEnvVar("WRTAG_URL"), "://")
	proto := elems[0]
	hostname := elems[1]

	wrUrl = fmt.Sprintf(wrPattern, proto, wrToken, hostname, wrOp)

	// Validate wrtag access
	code, _ := util.Request("GET", fmt.Sprintf("%s://%s@%s", proto, wrToken, hostname), strings.NewReader(""), "application/x-www-form-urlencoded")
	if code != 200 {
		log.Fatal("[ERROR] Cannot authenticate with wrtag, check WRTAG_TOKEN is configured correctly")
		os.Exit(1)
	}

	log.Println("[INFO] Loaded wrtag environment variables")
}
