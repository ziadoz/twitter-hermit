package saver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/dghubble/go-twitter/twitter"
)

type TweetSaver struct {
	SaveDir   string
	SaveJson  bool
	SaveMedia bool
	SaveLinks bool
}

func (ts *TweetSaver) Save(tweets []twitter.Tweet) error {
	var eg errgroup.Group
	var wg sync.WaitGroup
	links := make(chan string)

	for _, tweet := range tweets {
		if ts.SaveJson {
			tweetDir := path.Join(ts.SaveDir, strconv.FormatInt(tweet.ID, 10))
			if err := makeDir(tweetDir); err != nil {
				return err
			}

			dest := tweetDir
			twt := tweet
			eg.Go(func() error {
				return ts.saveJson(dest, twt)
			})
		}

		if ts.SaveMedia && hasMedia(tweet) {
			tweetDir := path.Join(ts.SaveDir, strconv.FormatInt(tweet.ID, 10))
			if err := makeDir(tweetDir); err != nil {
				return err
			}

			dest := tweetDir
			twt := tweet
			eg.Go(func() error {
				return ts.saveMedia(dest, twt)
			})
		}

		if ts.SaveLinks && hasLinks(tweet) {
			wg.Add(1)
			go ts.saveLinks(links, &wg, tweet)
		}
	}

	if ts.SaveJson || ts.SaveMedia {
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	if ts.SaveLinks {
		go func() {
			wg.Wait()
			close(links)
		}()

		file, err := makeLinks(path.Join(ts.SaveDir, "links.txt"))
		if err != nil {
			return err
		}
		defer file.Close()

		for link := range links {
			file.WriteString(link + "\n")
		}
	}

	return nil
}

func (ts *TweetSaver) saveJson(dest string, tweet twitter.Tweet) error {
	bytes, err := json.MarshalIndent(tweet, "", "    ")
	if err != nil {
		return fmt.Errorf("could not marshal tweet JSON: %s", err)
	}

	if err := ioutil.WriteFile(path.Join(dest, "tweet.json"), bytes, 0744); err != nil {
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

func (ts *TweetSaver) saveLinks(links chan<- string, wg *sync.WaitGroup, tweet twitter.Tweet) {
	defer wg.Done()
	for _, url := range tweet.Entities.Urls {
		links <- url.ExpandedURL
	}
}

func hasLinks(tweet twitter.Tweet) bool {
	return tweet.Entities != nil && tweet.Entities.Urls != nil && len(tweet.Entities.Urls) > 0
}

func makeDir(dest string) error {
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(dest, 0744); err != nil {
		return fmt.Errorf("could not make output directory: %s", err)
	}

	return nil
}

func makeLinks(dest string) (*os.File, error) {
	file, err := os.OpenFile(dest, os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_WRONLY, 0744)
	if err != nil {
		return nil, fmt.Errorf("could not create links file: %s", err)
	}

	return file, nil
}
