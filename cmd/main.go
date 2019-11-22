package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/karrick/tparse"
	"github.com/ziadoz/twitter-hermit/pkg/data"
	"github.com/ziadoz/twitter-hermit/pkg/hermit"
	"github.com/ziadoz/twitter-hermit/pkg/pathflag"
	"github.com/ziadoz/twitter-hermit/pkg/saver"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
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

	consumerKey := getRequiredEnv("TWITTER_CONSUMER_KEY")
	consumerSecret := getRequiredEnv("TWITTER_CONSUMER_SECRET")
	accessToken := getRequiredEnv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := getRequiredEnv("TWITTER_ACCESS_TOKEN_SECRET")

	flag.StringVar(&maxAge, "max-age", "-1month", "The max age of tweets to keep (e.g. -1day, -2weeks, -3months, -4years)")
	flag.Var(&saveDir, "save-dir", "Directory to save tweet content to")
	flag.BoolVar(&saveJson, "save-json", true, "Save tweet JSON?")
	flag.BoolVar(&saveMedia, "save-media", true, "Save tweet media?")
	flag.BoolVar(&dryRun, "dry-run", false, "Perform a dry run")
	flag.BoolVar(&silent, "silent", false, "Silence all log summary output")
	flag.Parse()

	if maxAge == "" {
		log.Fatal("missing max-age argument")
	}

	if strings.IndexRune(maxAge, '-') == -1 {
		log.Fatal("max-age argument must be negative")
	}

	now := time.Now()
	maxAgeTime, err := tparse.AddDuration(now, maxAge)
	if err != nil {
		log.Fatalf("invalid max age argument: %s\n", err)
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

	client := getTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	destroyer := &hermit.Destroyer{
		MaxAge:     maxAgeTime,
		DryRun:     dryRun,
		Output:     logger,
		TweetSaver: saver,
	}

	fmt.Fprintln(logger, "Twitter Hermit")
	fmt.Fprintln(logger, "==============")
	fmt.Fprintln(logger, "Max Age: ", maxAgeTime.Format("2 Jan 2006 03:04pm"))

	tweetErr := destroyer.Destroy(&data.UserTweets{Twitter: client})
	if tweetErr != nil {
		log.Fatal(tweetErr)
	}

	favouriteErr := destroyer.Destroy(&data.UserFavourites{Twitter: client})
	if favouriteErr != nil {
		log.Fatal(favouriteErr)
	}
}

// Get a required environment variable or panic.
// https://blog.antoine-augusti.fr/2015/12/testing-an-os-exit-scenario-in-golang/
func getRequiredEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Missing required environment variable %s\n", name)
	}
	return val
}

// Get a configured twitter.Client.
func getTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *twitter.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	http := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(http)
}
