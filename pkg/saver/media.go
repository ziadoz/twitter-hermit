package saver

import (
	"errors"
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
func extractMedia(media []twitter.MediaEntity) []string {
	links := make([]string, len(media))
	for i, source := range media {
		switch source.Type {
		case "photo":
			links[i] = source.MediaURLHttps

		case "animated_gif", "video":
			links[i] = source.VideoInfo.Variants[0].URL
		}
	}
	return links
}

// Determine the extension from a URL.
func getExtensionFromURL(link string) (string, error) {
	parts, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("could not parse URL: %s", err)
	}

	ext := path.Ext(parts.Path)
	if ext == "" {
		return "", errors.New("could not get extension from URL")
	}

	return ext, nil
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
