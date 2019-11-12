package saver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	dir := path.Join(ts.SaveDir, tweetId)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0744); err != nil {
			return fmt.Errorf("could not create tweet directory: %s", err)
		}
	}

	if ts.SaveJson {
		bytes, err := json.MarshalIndent(tweet, "", "    ")
		if err != nil {
			return fmt.Errorf("could not marshal tweet JSON: %s", err)
		}

		if err := ioutil.WriteFile(path.Join(dir, "tweet.json"), bytes, 0644); err != nil {
			return fmt.Errorf("could not write JSON file: %s", err)
		}
	}

	if ts.SaveMedia {
		saveMedia(dir, extractMedia(tweet))
	}

	return nil
}
