package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/ziadoz/twitter-hermit/pkg/hermit"
	"github.com/ziadoz/twitter-hermit/pkg/twitter"
	"github.com/ziadoz/twitter-hermit/pkg/util"
)

// Development Links
// https://developer.twitter.com/en/docs/api-reference-index
// https://developer.twitter.com/en/docs/tweets/timelines/guides/working-with-timelines.html
func main() {
	log.SetFlags(0)

	consumerKey := util.GetRequiredEnv("TWITTER_CONSUMER_KEY")
	consumerSecret := util.GetRequiredEnv("TWITTER_CONSUMER_SECRET")
	accessToken := util.GetRequiredEnv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := util.GetRequiredEnv("TWITTER_ACCESS_TOKEN_SECRET")

	var maxAge string
	var dryRun bool
	var silent bool
	var extractLinks string
	flag.StringVar(&maxAge, "max-age", "1 month", "The max age tweets to keep (e.g. 1 day, 2 weeks, 3 months, 4 years)")
	flag.BoolVar(&dryRun, "dry-run", false, "Performs a dry run that only outputs a log summary.")
	flag.BoolVar(&silent, "silent", false, "Silences all log summary output.")
	flag.StringVar(&extractLinks, "extract-links", "", "A text file to extract links from deleted tweets to.")
	flag.Parse()

	if maxAge == "" {
		log.Fatal("missing max-age flag")
	}

	maxAgeTime, err := util.ParseMaxAge(maxAge)
	if err != nil {
		log.Fatalf("invalid max age flag: %s\n", err)
	}

	var logger io.Writer
	logger = os.Stdout
	if !dryRun && silent {
		logger = ioutil.Discard
	}

	var linksFile io.Writer
	if extractLinks != "" {
		linksFile, err = os.OpenFile(extractLinks, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("invalid extract link file: %s", err)
		}
	}

	client := util.GetTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	destroyer := &hermit.Destroyer{
		MaxAge: maxAgeTime,
		DryRun: dryRun,
		Output: logger,
		Links:  linksFile,
	}

	tweetErr := destroyer.Destroy(&twitter.UserTweets{Twitter: client})
	if tweetErr != nil {
		log.Fatal(tweetErr)
	}
	fmt.Println()

	favouriteErr := destroyer.Destroy(&twitter.UserFavourites{Twitter: client})
	if favouriteErr != nil {
		log.Fatal(favouriteErr)
	}
	fmt.Println()
}
