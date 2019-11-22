package saver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

// Paths to test files.
var (
	fixture = "./fixtures/omg.gif"
	output  = "./output/omg.gif"
)

// Clean up any old files.
func init() {
	os.Remove(output)
}

func TestGetExtensionFromURL(t *testing.T) {
	is := is.New(t)
	ext, err := getExtensionFromURL("http://www.example.com/path/to/1234567890_foo_bar_bar.jpg")
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

	// Test web server that serves up fixture.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fixture)
	}))
	defer ts.Close()

	// Compare byte sizes.
	_, err := saveMediaFromURL(ts.URL, output)
	is.NoErr(err)

	// Compare exact number of bytes.
	fbytes, err := ioutil.ReadFile(fixture)
	obytes, err := ioutil.ReadFile(output)
	is.Equal(obytes, fbytes)
}
