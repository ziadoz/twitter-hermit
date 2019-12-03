package data

import (
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/matryer/is"
)

var tweets = []twitter.Tweet{
	{
		ID:        1,
		CreatedAt: "Tue Jan 01 01:01:01 -0000 2019",
	},
	{
		ID:        2,
		CreatedAt: "Wed Jan 02 02:02:02 -0000 2019",
	},
	{
		ID:        3,
		CreatedAt: "Thu Jan 03 03:03:03 -0000 2019",
	},
}

func TestFilterTweets(t *testing.T) {
	is := is.New(t)
	maxAge, err := time.Parse(time.RubyDate, "Wed Jan 02 02:02:02 -0000 2019")
	is.NoErr(err)
	is.Equal([]twitter.Tweet{tweets[0]}, FilterTweets(tweets, maxAge))
}

func TextGetMaxID(t *testing.T) {
	is := is.New(t)
	is.Equal(3, GetMaxID(tweets))
}
