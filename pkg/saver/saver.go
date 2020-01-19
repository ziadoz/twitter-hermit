package saver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

type TweetSaver struct {
	SaveDir   string
	SaveJson  bool
	SaveMedia bool
	SaveLinks bool
}

func (ts *TweetSaver) Save(tweets []twitter.Tweet) error {
	var linksFile *os.File
	defer linksFile.Close()

	if ts.SaveLinks {
		file, err := makeLinks(path.Join(ts.SaveDir, "links.txt"))
		if err != nil {
			return err
		}

		linksFile = file
	}

	for _, tweet := range tweets {
		if ts.SaveJson {
			tweetDir := path.Join(ts.SaveDir, strconv.FormatInt(tweet.ID, 10))
			if err := makeDir(tweetDir); err != nil {
				return err
			}

			if err := ts.saveJson(tweetDir, tweet); err != nil {
				return err
			}
		}

		if ts.SaveMedia && hasMedia(tweet) {
			tweetDir := path.Join(ts.SaveDir, strconv.FormatInt(tweet.ID, 10))
			if err := makeDir(tweetDir); err != nil {
				return err
			}

			if err := ts.saveMedia(tweetDir, tweet); err != nil {
				return err
			}
		}

		if ts.SaveLinks && hasLinks(tweet) {
			if err := ts.saveLinks(linksFile, tweet); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ts *TweetSaver) saveJson(dest string, tweet twitter.Tweet) error {
	bytes, err := json.MarshalIndent(tweet, "", "    ")
	if err != nil {
		return fmt.Errorf("could not marshal tweet JSON: %s", err)
	}

	if err := ioutil.WriteFile(path.Join(dest, "tweet.json"), bytes, 0755); err != nil {
		return fmt.Errorf("could not write JSON file: %s", err)
	}

	return nil
}

func (ts *TweetSaver) saveMedia(dest string, tweet twitter.Tweet) error {
	num := 1
	for _, media := range extractMedia(tweet.ExtendedEntities.Media) {
		ext, err := getExtensionFromURL(media)
		if err != nil {
			return fmt.Errorf("could not save tweet media: %s", media)
		}

		saveMediaFromURL(media, path.Join(dest, "media-"+strconv.Itoa(num)+ext))
		num++
	}

	return nil
}

func hasMedia(tweet twitter.Tweet) bool {
	return tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0
}

func (ts *TweetSaver) saveLinks(file *os.File, tweet twitter.Tweet) error {
	links := make([]string, len(tweet.Entities.Urls))
	for i, url := range tweet.Entities.Urls {
		links[i] = url.ExpandedURL
	}

	if _, err := file.WriteString(strings.Join(links, "\n") + "\n"); err != nil {
		return fmt.Errorf("could not save links: %s", err)
	}

	return nil
}

func hasLinks(tweet twitter.Tweet) bool {
	return tweet.Entities != nil && tweet.Entities.Urls != nil && len(tweet.Entities.Urls) > 0
}

func makeDir(dest string) error {
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("could not make output directory: %s", err)
	}

	return nil
}

func makeLinks(dest string) (*os.File, error) {
	file, err := os.OpenFile(dest, os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return nil, fmt.Errorf("could not create links file: %s", err)
	}

	return file, nil
}
