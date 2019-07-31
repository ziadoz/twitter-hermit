package twitter

import (
	"fmt"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

// Repository defines a way to get and a destroy a particular Twitter entity (e.g. tweets, links, favourites).
type Repository interface {
	Description() string
	Get(params QueryParams) ([]twitter.Tweet, error)
	Destroy(tweets []twitter.Tweet) error
}

type UserTweets struct {
	Twitter *twitter.Client
}

type QueryParams struct {
	Count int
	MaxID int64
}

func (t *UserTweets) Description() string {
	return "tweets"
}

func (t *UserTweets) Get(params QueryParams) ([]twitter.Tweet, error) {
	tweets, _, err := t.Twitter.Timelines.UserTimeline(&twitter.UserTimelineParams{
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

func (t *UserTweets) Destroy(tweets []twitter.Tweet) error {
	for _, tweet := range tweets {
		_, _, err := t.Twitter.Statuses.Destroy(tweet.ID, &twitter.StatusDestroyParams{
			TrimUser: twitter.Bool(true),
		})

		if err != nil {
			return fmt.Errorf("failed to delete tweet %s", err)
		}
	}

	return nil
}

type UserFavourites struct {
	Twitter *twitter.Client
}

func (f *UserFavourites) Description() string {
	return "favourites"
}

func (f *UserFavourites) Get(params QueryParams) ([]twitter.Tweet, error) {
	favorites, _, err := f.Twitter.Favorites.List(&twitter.FavoriteListParams{
		Count: params.Count,
		MaxID: params.MaxID,
	})

	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (f *UserFavourites) Destroy(favourites []twitter.Tweet) error {
	for _, favourite := range favourites {
		_, _, err := f.Twitter.Favorites.Destroy(&twitter.FavoriteDestroyParams{
			ID: favourite.ID,
		})

		if err != nil {
			log.Printf("Error: Failed to delete favourite %s", err)
		}
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
