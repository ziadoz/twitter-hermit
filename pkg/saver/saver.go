package saver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
)

type TweetSaver struct {
	SaveDir   string
	SaveJson  bool
	SaveMedia bool
}

func (ts *TweetSaver) Save(tweet twitter.Tweet) error {
	tweetId := strconv.FormatInt(tweet.ID, 10)

	if ts.SaveJson {
		bytes, err := json.MarshalIndent(tweet, "", "    ")
		if err != nil {
			return fmt.Errorf("could not marshal tweet JSON: %s", err)
		}

		if err := ioutil.WriteFile(path.Join(ts.SaveDir, tweetId+".json"), bytes, 0644); err != nil {
			return fmt.Errorf("could not write JSON file: %s", err)
		}
	}

	// @todo: Save files in concurrency.
	if ts.SaveMedia {
		num := 1
		for _, media := range extractMedia(tweet) {
			ext, err := getExtensionFromURL(media)
			if err != nil {
				return fmt.Errorf("could not save tweet ID %s media: %s", tweetId, media)
			}

			// Possibly an m3u8 file
			// https://gerardnico.com/video/m3u
			if ext == "" {
				continue
			}

			saveMediaFromURL(media, path.Join(ts.SaveDir, tweetId+"-"+strconv.Itoa(num)+ext))
			num++
		}
	}

	return nil
}
