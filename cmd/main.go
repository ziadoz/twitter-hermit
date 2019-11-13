package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/ziadoz/twitter-hermit/pkg/hermit"
	"github.com/ziadoz/twitter-hermit/pkg/pathflag"
	"github.com/ziadoz/twitter-hermit/pkg/saver"
	"github.com/ziadoz/twitter-hermit/pkg/twitter"
	"github.com/ziadoz/twitter-hermit/pkg/util"
)

// Development Links
// https://developer.twitter.com/en/docs/api-reference-index
// https://developer.twitter.com/en/docs/tweets/timelines/guides/working-with-timelines.html
func main() {
	run()
}

var (
	maxAge    string
	dryRun    bool
	silent    bool
	saveDir   pathflag.Path
	saveJson  bool
	saveMedia bool
)

func run() {
	log.SetFlags(0)

	consumerKey := util.GetRequiredEnv("TWITTER_CONSUMER_KEY")
	consumerSecret := util.GetRequiredEnv("TWITTER_CONSUMER_SECRET")
	accessToken := util.GetRequiredEnv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := util.GetRequiredEnv("TWITTER_ACCESS_TOKEN_SECRET")

	flag.StringVar(&maxAge, "max-age", "1 month", "The max age tweets to keep (e.g. 1 day, 2 weeks, 3 months, 4 years)")
	flag.Var(&saveDir, "save-dir", "Directory to save tweet content to")
	flag.BoolVar(&saveJson, "save-json", true, "Save tweet JSON?")
	flag.BoolVar(&saveMedia, "save-media", true, "Save tweet media?")
	flag.BoolVar(&dryRun, "dry-run", false, "Perform a dry run that only outputs a log summary")
	flag.BoolVar(&silent, "silent", false, "Silence all log summary output")
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

	saver := &saver.TweetSaver{
		SaveDir:   saveDir.Path,
		SaveJson:  saveJson,
		SaveMedia: saveMedia,
	}

	client := util.GetTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	destroyer := &hermit.Destroyer{
		MaxAge:     maxAgeTime,
		DryRun:     dryRun,
		Output:     logger,
		TweetSaver: saver,
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
