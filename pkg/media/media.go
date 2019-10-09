package media

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
)

type Media struct {
	Src  string
	Dest string
}

// Extract media from tweets.
// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/extended-entities-object
func Extract(tweets []twitter.Tweet) []Media {
	medias := []Media{}
	pattern := regexp.MustCompile(`\?.*$`) // Remove ?tag=123 query strings from URL. Should use url.Parse() in the future.

	for _, tweet := range tweets {
		if tweet.ExtendedEntities == nil || len(tweet.ExtendedEntities.Media) == 0 {
			continue
		}

		for _, source := range tweet.ExtendedEntities.Media {
			media := Media{}

			switch source.Type {
			case "photo":
				media.Src = source.MediaURLHttps
				media.Dest = pattern.ReplaceAllString(path.Base(source.MediaURLHttps), "")

			case "animated_gif":
			case "video":
				media.Src = source.VideoInfo.Variants[0].URL
				media.Dest = pattern.ReplaceAllString(path.Base(source.VideoInfo.Variants[0].URL), "")
			}

			medias = append(medias, media)
		}
	}

	return medias
}

// Save media to a directory.
func Save(dir string, medias []Media) error {
	sem := make(chan struct{}, 20)
	var wg sync.WaitGroup

	for _, media := range medias {
		wg.Add(1)
		go func(dir string, media Media) {
			defer wg.Done()
			sem <- struct{}{}
			_, err := SaveMediaFromURL(media.Src, path.Join(dir, media.Dest)) // handle errors somehow
			if err != nil {
				fmt.Println(err)
			}
			<-sem
		}(dir, media)
	}

	wg.Wait()
	return nil
}

// Save media file from a URL.
func SaveMediaFromURL(src, dest string) (int64, error) {
	file, err := os.Create(dest)
	if err != nil {
		return 0, fmt.Errorf("could not create file: %s", err)
	}
	defer file.Close()

	resp, err := http.Get(src)
	if err != nil {
		return 0, fmt.Errorf("could not request image: %s", err)
	}
	defer resp.Body.Close()

	bytes, err := io.Copy(file, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("could not copy image: %s", err)
	}

	return bytes, nil
}
