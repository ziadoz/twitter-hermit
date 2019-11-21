package util

import (
	"log"
	"os"
	"regexp"
	"strings"
)

// Get a required environment variable or panic.
// https://blog.antoine-augusti.fr/2015/12/testing-an-os-exit-scenario-in-golang/
func GetRequiredEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Missing required environment variable %s\n", name)
	}
	return val
}

// Strip any newlines from a string.
func StripNewlines(s string) string {
	re := regexp.MustCompile(`(\r|\n|\r\n)+`)
	return strings.Trim(re.ReplaceAllString(s, " "), " ")
}
