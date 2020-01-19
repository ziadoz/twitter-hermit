package saver

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/matryer/is"
)

var tweetId int64 = 1234567890

func init() {
	files := []string{
		"./output/1234567890/tweet.json",
		"./output/1234567890/media-1.gif",
		"./output/links.txt",
		"./output/omg.gif",
	}

	for _, file := range files {
		os.Remove(file)
	}
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

	err := ts.Save([]twitter.Tweet{tweet})
	is.NoErr(err)

	fbytes, _ := ioutil.ReadFile("./fixtures/1234567890/tweet.json")
	obytes, _ := ioutil.ReadFile("./output/1234567890/tweet.json")
	is.True(bytes.Equal(fbytes, obytes))
}

func TestTweetSaverSaveMedia(t *testing.T) {
	is := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./fixtures/omg.gif")
	}))
	defer ts.Close()

	url := ts.URL + "/omg.gif" // Add filename and extension to URL.
	tweet := twitter.Tweet{
		ID:        tweetId,
		CreatedAt: "2019-01-01 01:01:01",
		Text:      "This is a test tweet. @helloworld http://www.example.com",
		ExtendedEntities: &twitter.ExtendedEntity{
			Media: []twitter.MediaEntity{
				twitter.MediaEntity{
					MediaURLHttps: url,
					Type:          "animated_gif",
					VideoInfo: twitter.VideoInfo{
						Variants: []twitter.VideoVariant{
							twitter.VideoVariant{
								URL: url,
							},
						},
					},
				},
			},
		},
	}

	saver := TweetSaver{
		SaveDir:   "./output",
		SaveMedia: true,
	}

	err := saver.Save([]twitter.Tweet{tweet})
	is.NoErr(err)

	fbytes, _ := ioutil.ReadFile("./fixtures/1234567890/media-1.gif")
	obytes, _ := ioutil.ReadFile("./output/1234567890/media-1.gif")
	is.True(bytes.Equal(fbytes, obytes))
}

func TestTweetSaverSaveLinks(t *testing.T) {
	is := is.New(t)

	tweet := twitter.Tweet{
		ID: tweetId,
		Entities: &twitter.Entities{
			Urls: []twitter.URLEntity{
				{
					ExpandedURL: "http://www.example1.com",
				},
				{
					ExpandedURL: "http://www.example2.com",
				},
			},
		},
	}

	saver := TweetSaver{
		SaveDir:   "./output",
		SaveLinks: true,
	}

	err := saver.Save([]twitter.Tweet{tweet})
	is.NoErr(err)

	fbytes, _ := ioutil.ReadFile("./fixtures/links.txt")
	obytes, _ := ioutil.ReadFile("./output/links.txt")
	is.True(bytes.Equal(fbytes, obytes))
}
