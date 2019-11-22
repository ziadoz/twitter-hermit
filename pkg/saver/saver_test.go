package saver

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/matryer/is"
)

var tweetId int64 = 1234567890

func init() {
	os.Remove("./output/1234567890.json")
}

func TestTweetSaverSaveJson(t *testing.T) {
	is := is.New(t)

	tweet := twitter.Tweet{
		ID:        tweetId,
		CreatedAt: "2019-01-01 01:01:01",
		Text:      "This is a test tweet. @helloworld http://www.example.com",
	}

	ts := TweetSaver{
		SaveDir:  "./output",
		SaveJson: true,
	}

	err := ts.Save(tweet)
	is.NoErr(err)

	fbytes, _ := ioutil.ReadFile("./fixtures/1234567890.json")
	obytes, _ := ioutil.ReadFile("./output/1234567890.json")
	is.True(bytes.Equal(fbytes, obytes))
}
