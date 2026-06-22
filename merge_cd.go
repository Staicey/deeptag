package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/staicey/deeptag/deepcrate"
)

// Merge releases into a single directory
func mergeReleaseCDs(event deepcrate.WebHook) string {
	// Only merge if the download dir matches a known pattern
	if !cdRegex.MatchString(event.Data.DownloadPath) {
		log.Println("[DEBUG] mergeReleaseCDs: Does not appear to be a CD release")
		return event.Data.DownloadPath
	}

	// Make a new path to store the combined CDs in
	newPath := filepath.Join(dlDir, stripRegex.ReplaceAllString(event.Data.Album, ""))
	if !strings.HasPrefix(newPath, dlDir) {
		log.Println("[ERROR] Malicious album path detected, skipping merge")
		log.Println("[DEBUG] mergeReleaseCDs:", newPath)
		return event.Data.DownloadPath
	}

	// Create the new directory
	os.Mkdir(newPath, 0750)

	// Glob the download directory for disks matching the pattern
	globPattern := filepath.Join(dlDir, alphaRegex.ReplaceAllString(event.Data.DownloadPath, ""))
	cds, _ := filepath.Glob(globPattern + "*")
	log.Println("[DEBUG] mergeReleaseCDs: globbed dirs [", globPattern+"*]:", cds)
	for _, cd := range cds {
		files, _ := filepath.Glob(filepath.Join(cd, "*"))
		log.Println("[DEBUG] mergeReleaseCDs: globbed files[", filepath.Join(cd, "*")+"]:", files)

		for _, f := range files {
			log.Printf("[DEBUG] Moving %s -> %s\n", f, filepath.Join(newPath, fmt.Sprintf("%s-%s", filepath.Base(cd), filepath.Base(f))))
			if err := os.Rename(f, filepath.Join(newPath, fmt.Sprintf("%s-%s", filepath.Base(cd), filepath.Base(f)))); err != nil {
				log.Println("[ERROR]", err)
			}
		}

		if err := os.Remove(cd); err != nil {
			log.Println("[ERROR]", err)
		}
	}

	return newPath
}
