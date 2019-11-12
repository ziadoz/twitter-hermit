package saver

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

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

// Save media links to a directory.
// @todo: REFACTOR GOROUTINES TO WORK PROPERLY.
func saveMedia(dir string, links []string) {
	sem := make(chan struct{}, 20)
	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go func(dir, link string) {
			defer wg.Done()
			sem <- struct{}{}
			dest, err := getFileNameFromURL(link)
			if err != nil {
				// @todo: HANDLE ERROR.
			}

			_, err = saveMediaFromURL(link, path.Join(dir, dest))
			if err != nil {
				// @todo: HANDLE ERROR.
			}
			<-sem
		}(dir, link)
	}

	wg.Wait()
	close(sem)
}

// Determine the filename from a URL.
func getFileNameFromURL(src string) (string, error) {
	parts, err := url.Parse(src)
	if err != nil {
		return "", fmt.Errorf("could not get filename from URL: %s", err)
	}

	return path.Base(parts.Path), nil
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
