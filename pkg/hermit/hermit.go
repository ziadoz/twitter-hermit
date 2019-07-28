package hermit

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/ziadoz/twitter-hermit/pkg/util"
)

var logFormat string = " - ID:   %d\n   Date: %s\n   Text: %s\n"

type Client struct {
	Twitter *twitter.Client
	Writer  io.Writer
	DryRun  bool
}

type QueryParams struct {
	Count int
	MaxID int64
}

func (c *Client) GetUserTweets(params QueryParams) ([]twitter.Tweet, error) {
	tweets, _, err := c.Twitter.Timelines.UserTimeline(&twitter.UserTimelineParams{
		Count:           params.Count,
		MaxID:           params.MaxID,
		IncludeRetweets: twitter.Bool(true),
		TrimUser:        twitter.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func (c *Client) DestroyTweets(tweets []twitter.Tweet) error {
	for _, tweet := range tweets {
		if !c.DryRun {
			_, _, err := c.Twitter.Statuses.Destroy(tweet.ID, &twitter.StatusDestroyParams{
				TrimUser: twitter.Bool(true),
			})

			if err != nil {
				return fmt.Errorf("Failed to delete tweet %s", err)
			}
		}

		createdAt, _ := tweet.CreatedAtTime()
		fmt.Fprintf(
			c.Writer,
			logFormat,
			tweet.ID,
			createdAt.Format("2 Jan 2006 03:04pm"),
			util.StripNewlines(tweet.Text),
		)
	}

	return nil
}

func (c *Client) GetUserFavourites(params QueryParams) ([]twitter.Tweet, error) {
	favorites, _, err := c.Twitter.Favorites.List(&twitter.FavoriteListParams{
		Count: params.Count,
		MaxID: params.MaxID,
	})

	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (c *Client) DestroyFavourites(favourites []twitter.Tweet) error {
	for _, favourite := range favourites {
		if !c.DryRun {
			_, _, err := c.Twitter.Favorites.Destroy(&twitter.FavoriteDestroyParams{
				ID: favourite.ID,
			})

			if err != nil {
				log.Printf("Error: Failed to delete favourite %s", err)
			}
		}

		createdAt, _ := favourite.CreatedAtTime()
		fmt.Fprintf(
			c.Writer,
			logFormat,
			favourite.ID,
			createdAt.Format("2 Jan 2006 03:04pm"),
			util.StripNewlines(favourite.Text),
		)
	}

	return nil
}

// Filter out tweets that are newer than the max age time.Time.
func FilterTweets(tweets []twitter.Tweet, maxAge time.Time) []twitter.Tweet {
	var filtered = make([]twitter.Tweet, 0)

	for _, tweet := range tweets {
		createdAt, _ := tweet.CreatedAtTime()
		if createdAt.Before(maxAge) {
			filtered = append(filtered, tweet)
		}
	}

	return filtered
}

// Return the MaxID (last tweet's ID) from a slice of tweets.
func GetMaxID(tweets []twitter.Tweet) int64 {
	return tweets[len(tweets)-1].ID
}
