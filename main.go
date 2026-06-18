package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/logutils"

	"github.com/staicey/deeptag/deepcrate"
	"github.com/staicey/deeptag/wrtag"
)

func main() {
	// Setup logger
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"FINE", "DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stdout,
	}

	if _, ok := os.LookupEnv("DEBUG"); ok == true {
		filter.MinLevel = "DEBUG"
	}

	log.Default().SetFlags(log.Ltime | log.Ldate | log.LUTC)
	log.Default().SetOutput(filter)

	// Register webhook route
	r := mux.NewRouter()
	r.HandleFunc("/import", TransformRequest).Methods("POST")

	log.Println("[INFO] Listening at port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

// Converts a webhook from deepcrate -> wrtag import
func TransformRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the webhook from deepcrate
	var event deepcrate.WebHook
	json.NewDecoder(r.Body).Decode(&event)

	log.Printf("[INFO] Request from %s %#v\n", r.RemoteAddr, event)
	// If this is a test event from deepcrate then don't take actions
	if event.Data.Album == "Test Album" && event.Data.Artist == "Test Artist" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// If the download path is empty then there was a failure downloading from
	// slskd so we can remove the download entry & re-add to wishlist
	if len(event.Data.DownloadPath) == 0 {
		id := deepcrate.GetDownloadID(&event)
		deepcrate.DeleteDownloadRecord(id)
		time.Sleep(time.Second * 5)

		deepcrate.AddToWishlist(&event)
		time.Sleep(time.Second * 2)
		deepcrate.Trigger("slskd-downloader")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Trigger an import via wrtag & return the result to deepcrate
	code, output := wrtag.Import(event.Data.DownloadPath)
	w.WriteHeader(code)
	w.Write(output)
}
