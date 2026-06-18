package util

import (
	"log"
	"os"
)

func LoadEnvVar(key string) string {
	if val, ok := os.LookupEnv(key); ok == true {
		return val
	}

	log.Fatal("[ERROR] Failed to load environment variables:", key)
	os.Exit(1)
	return ""
}
