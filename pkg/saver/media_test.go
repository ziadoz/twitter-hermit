package saver

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/matryer/is"
)

var (
	fixture string = "./fixtures/omg.gif"
	output  string = "./output/omg.gif"
)

func init() {
	os.Remove(output)
}

func TextExtractMedia(t *testing.T) {
	is := is.New(t)

	have := []twitter.MediaEntity{
		twitter.MediaEntity{
			MediaURLHttps: "https://www.example.com/path/to/1234567890_foo_bar_baz.jpg",
			Type:          "photo",
		},
		twitter.MediaEntity{
			MediaURLHttps: "https://www.example.com/path/to/1234567890_foo_bar_baz.mp4",
			Type:          "video",
			VideoInfo: twitter.VideoInfo{
				Variants: []twitter.VideoVariant{
					twitter.VideoVariant{
						URL: "https://www.example.com/path/to/1234567890_foo_bar_baz_large.mp4",
					},
					twitter.VideoVariant{
						URL: "https://www.example.com/path/to/1234567890_foo_bar_baz_medium.mp4",
					},
					twitter.VideoVariant{
						URL: "https://www.example.com/path/to/1234567890_foo_bar_baz_small.mp4",
					},
				},
			},
		},
		twitter.MediaEntity{
			MediaURLHttps: "https://www.example.com/path/to/1234567890_foo_bar_baz.gif",
			Type:          "animated_gif",
			VideoInfo: twitter.VideoInfo{
				Variants: []twitter.VideoVariant{
					twitter.VideoVariant{
						URL: "https://www.example.com/path/to/1234567890_foo_bar_baz_large.gif",
					},
					twitter.VideoVariant{
						URL: "https://www.example.com/path/to/1234567890_foo_bar_baz_medium.gif",
					},
					twitter.VideoVariant{
						URL: "https://www.example.com/path/to/1234567890_foo_bar_baz_small.gif",
					},
				},
			},
		},
	}

	want := []string{
		"https://www.example.com/path/to/1234567890_foo_bar_baz.jpg",
		"https://www.example.com/path/to/1234567890_foo_bar_baz_large.mp4",
		"https://www.example.com/path/to/1234567890_foo_bar_baz_large.gif",
	}

	is.Equal(extractMedia(have), want)
}

func TestGetExtensionFromURL(t *testing.T) {
	is := is.New(t)
	ext, err := getExtensionFromURL("http://www.example.com/path/to/1234567890_foo_bar_baz.jpg")
	is.Equal(ext, ".jpg")
	is.NoErr(err)
}

func TestGetExtensionFromURLNoExtension(t *testing.T) {
	is := is.New(t)
	ext, err := getExtensionFromURL("http://www.example.com/path/to/nothing")
	is.Equal(ext, "")
	is.True(strings.Contains(err.Error(), "could not get extension from URL"))
}

func TestSaveMediaFromURL(t *testing.T) {
	is := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fixture)
	}))
	defer ts.Close()

	_, err := saveMediaFromURL(ts.URL, output)
	is.NoErr(err)

	fbytes, _ := ioutil.ReadFile(fixture)
	obytes, _ := ioutil.ReadFile(output)
	is.True(bytes.Equal(obytes, fbytes))
}
