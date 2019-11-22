package saver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

	if ts.SaveMedia && hasMedia(tweet) {
		num := 1
		for _, media := range extractMedia(tweet.ExtendedEntities.Media) {
			log.Fatal(media)

			ext, err := getExtensionFromURL(media)
			if err != nil {
				return fmt.Errorf("could not save tweet ID %s media: %s", tweetId, media)
			}

			saveMediaFromURL(media, path.Join(ts.SaveDir, tweetId+"-"+strconv.Itoa(num)+ext))
			num++
		}
	}

	return nil
}

func hasMedia(tweet twitter.Tweet) bool {
	return (tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0)
}
