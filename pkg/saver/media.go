package saver

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/dghubble/go-twitter/twitter"
)

// Extract media links from tweets.
// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/extended-entities-object
func extractMedia(tweet twitter.Tweet) []string {
	if tweet.ExtendedEntities == nil || len(tweet.ExtendedEntities.Media) == 0 {
		return []string{}
	}

	links := make([]string, len(tweet.ExtendedEntities.Media))
	for _, source := range tweet.ExtendedEntities.Media {
		switch source.Type {
		case "photo":
			links = append(links, source.MediaURLHttps)

		case "animated_gif", "video":
			links = append(links, source.VideoInfo.Variants[0].URL)
		}
	}

	return links
}

// Determine the extension from a URL.
func getExtensionFromURL(link string) (string, error) {
	parts, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("could not get extension from URL: %s", err)
	}

	return path.Ext(parts.Path), nil
}

// Save media file from a URL.
func saveMediaFromURL(src, dest string) (int64, error) {
	file, err := os.Create(dest)
	if err != nil {
		return 0, fmt.Errorf("could not create file: %s", err)
	}
	defer file.Close()

	resp, err := http.Get(src)
	if err != nil {
		return 0, fmt.Errorf("could not request file: %s", err)
	}
	defer resp.Body.Close()

	bytes, err := io.Copy(file, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("could not copy file: %s", err)
	}

	return bytes, nil
}
