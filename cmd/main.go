package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ziadoz/twitter-hermit/pkg/hermit"
	"github.com/ziadoz/twitter-hermit/pkg/util"
)

const batchSize = 200

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

	var writer io.Writer
	writer = os.Stdout
	if !dryRun && silent {
		writer = ioutil.Discard
	}

	var linksFile io.Writer
	if extractLinks != "" {
		linksFile, err = os.OpenFile(extractLinks, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("invalid extract link file: %s", err)
		}
	}

	client := &hermit.Client{
		Twitter: util.GetTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret),
		Writer:  writer,
		DryRun:  dryRun,
	}

	fmt.Println("[Deleting Tweets]")
	var tweetMaxID int64
	var tweetDeletedCount int
	for {
		tweets, err := client.GetUserTweets(hermit.QueryParams{Count: batchSize, MaxID: tweetMaxID})
		if err != nil {
			log.Fatalf("could not get user tweets: %s\n", err)
		}

		if len(tweets) == 0 {
			break // We're done deleting.
		}

		filteredTweets := hermit.FilterTweets(tweets, maxAgeTime)
		if len(filteredTweets) == 0 {
			tweetMaxID = hermit.GetMaxID(tweets) - 1
			continue
		}

		if linksFile != nil {
			links := util.FollowLinkRedirects(util.ExtractLinks(filteredTweets))
			if len(links) > 0 {
				fmt.Fprintf(linksFile, strings.Join(links, "\n")+"\n")
			}
		}

		err = client.DestroyTweets(filteredTweets)
		if err != nil {
			log.Fatalf("could not delete tweets: %s\n", err)
		}

		tweetDeletedCount += len(filteredTweets)
		tweetMaxID = hermit.GetMaxID(tweets) - 1
	}

	if tweetDeletedCount > 0 {
		fmt.Printf("Deleted %d tweets successfully\n\n", tweetDeletedCount)
	} else {
		fmt.Printf("No tweets needed deleting\n\n")
	}

	fmt.Println("[Deleting Favourites]")
	var favouriteMaxID int64
	var favouriteDeletedCount int
	for {
		tweets, err := client.GetUserFavourites(hermit.QueryParams{Count: batchSize, MaxID: favouriteMaxID})
		if err != nil {
			log.Fatalf("could not get user tweets: %s\n", err)
		}

		if len(tweets) == 0 {
			break // We're done deleting.
		}

		filteredTweets := hermit.FilterTweets(tweets, maxAgeTime)
		if len(filteredTweets) == 0 {
			favouriteMaxID = hermit.GetMaxID(tweets) - 1
			continue
		}

		if linksFile != nil {
			links := util.FollowLinkRedirects(util.ExtractLinks(filteredTweets))
			if len(links) > 0 {
				fmt.Fprintf(linksFile, strings.Join(links, "\n")+"\n")
			}
		}

		err = client.DestroyFavourites(filteredTweets)
		if err != nil {
			log.Fatalf("could not delete favourites: %s\n", err)
		}

		favouriteDeletedCount += len(filteredTweets)
		favouriteMaxID = hermit.GetMaxID(tweets) - 1
	}

	if favouriteDeletedCount > 0 {
		fmt.Printf("Deleted %d favourites successfully\n\n", favouriteDeletedCount)
	} else {
		fmt.Printf("No tweets favourites deleting\n\n")
	}
}
