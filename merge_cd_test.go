package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/staicey/deeptag/deepcrate"
)

func doMergeCD(t *testing.T, cdPattern string) {
	// Create the file structure
	dlDir = "testing/test_download"
	for r := range 2 {
		os.MkdirAll(fmt.Sprintf("%s/%s%d", dlDir, cdPattern, r+1), 0750)
		for i := range 10 {
			os.Create(fmt.Sprintf("%s/%s%d/Track %d", dlDir, cdPattern, r+1, i+1))
		}
	}

	// Create a fake webhook event
	newPath := mergeReleaseCDs(deepcrate.WebHook{
		Event: "",
		Data: struct {
			Artist       string "json:\"artist\""
			Album        string "json:\"album\""
			DownloadPath string "json:\"download_path\""
		}{
			Artist:       "",
			Album:        "Test Album",
			DownloadPath: fmt.Sprintf("%s1", cdPattern),
		},
	})

	// Check the resulting filesystem looks correct
	assert.Equal(t, fmt.Sprintf("%s/Test Album", dlDir), newPath)
	for c := range 2 {
		for n := range 10 {
			assert.FileExists(t, fmt.Sprintf("%s/Test Album/%s%d-Track %d", dlDir, cdPattern, c+1, n+1))
		}
	}

	os.RemoveAll("./testing")
}

func TestDiskPattern(t *testing.T) {
	for _, pattern := range []string{"CD_", "CD", "Disk", "dIsK_"} {
		doMergeCD(t, pattern)
	}
}

func TestNonDiskPattern(t *testing.T) {
	for _, pattern := range []string{"Test Album", "2000 Album", "abcdef"} {
		newPath := mergeReleaseCDs(deepcrate.WebHook{Event: "", Data: struct {
			Artist       string "json:\"artist\""
			Album        string "json:\"album\""
			DownloadPath string "json:\"download_path\""
		}{DownloadPath: pattern}})

		assert.Equal(t, pattern, newPath)
	}
}
