package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Get a required environment variable or panic.
func GetRequiredEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Missing required environment variable %s\n", name)
	}
	return val
}

// Get a configured twitter.Client.
func GetTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *twitter.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	http := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(http)
}

// Parse a max age formatted string into a time.Time.
func ParseMaxAge(maxAge string) (time.Time, error) {
	pattern, _ := regexp.Compile(`(?P<length>\d+)\s+(?P<duration>day|week|month|year)s?`)
	matches := pattern.FindStringSubmatch(maxAge)
	if len(matches) == 0 {
		return time.Now(), fmt.Errorf("invalid duration %s", maxAge)
	}

	length, _ := strconv.Atoi(matches[1])
	duration := matches[2]
	if length == 0 {
		return time.Now(), errors.New("duration is zero")
	}

	years, months, days := 0, 0, 0
	switch duration {
	case "year":
		years = length * -1
	case "month":
		months = length * -1
	case "week":
		days = (7 * length) * -1
	case "day":
		days = length * -1
	}

	return time.Now().AddDate(years, months, days), nil
}

// Strip any newlines from a string.
func StripNewlines(s string) string {
	re := regexp.MustCompile(`(\r|\n|\r\n|\n\r)+`)
	return re.ReplaceAllString(s, " ")
}
