package hermit

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/matryer/is"
	"github.com/ziadoz/twitter-hermit/pkg/data"
)

var tweets = []twitter.Tweet{
	{
		ID:        1,
		CreatedAt: "Tue Jan 01 01:01:01 -0000 2019",
		Text:      "I am tweet number 1",
	},
	{
		ID:        2,
		CreatedAt: "Wed Jan 02 02:02:02 -0000 2019",
		Text:      "I am tweet number 2",
	},
	{
		ID:        3,
		CreatedAt: "Thu Jan 03 03:03:03 -0000 2019",
		Text:      "I am tweet number 3",
	},
	{
		ID:        4,
		CreatedAt: "Fri Jan 04 04:04:04 -0000 2019",
		Text:      "I am tweet number 4",
	},
	{
		ID:        5,
		CreatedAt: "Fri Jan 05 05:05:05 -0000 2019",
		Text:      "I am tweet number 5",
	},
}

type testrepo struct{
	GetErr error
	DestroyErr error
	position int
}

func (tr *testrepo) Description() string {
	return "tweets"
}

func (tr *testrepo) Get(params data.QueryParams) ([]twitter.Tweet, error) {
	if tr.GetErr != nil {
		return []twitter.Tweet{}, tr.GetErr
	}

	start := tr.position
	end := start + params.Count
	if start > len(tweets) || end > len(tweets) {
		return []twitter.Tweet{}, nil
	}

	tr.position += params.Count
	return tweets[start : end], nil
}

func (tr *testrepo) Destroy(tweets []twitter.Tweet) error {
	if tr.DestroyErr != nil {
		return tr.DestroyErr
	}

	return nil
}

func TestDestroy(t *testing.T) {
	is := is.New(t)
	var b bytes.Buffer
	maxAge, _ := time.Parse(time.RubyDate, "Thu Jan 03 03:03:03 -0000 2019")

	h := Destroyer{BatchSize: 2, MaxAge: maxAge, Output: &b}
	err := h.Destroy(&testrepo{})
	is.NoErr(err)
	is.Equal("• Deleted 2 tweets.\n", b.String())
}

func TestDestroyGetError(t *testing.T) {
	is := is.New(t)
	tr := &testrepo{GetErr: errors.New("something went wrong")}

	h := Destroyer{}
	err := h.Destroy(tr)
	is.Equal("could not get tweets: something went wrong\n", err.Error())
}

func TestDestroyDestroyError(t *testing.T) {
	is := is.New(t)
	tr := &testrepo{DestroyErr: errors.New("something went wrong")}
	maxAge, _ := time.Parse(time.RubyDate, "Thu Jan 03 03:03:03 -0000 2019")

	h := Destroyer{BatchSize: 2, MaxAge: maxAge}
	err := h.Destroy(tr)
	is.Equal("could not destroy tweets: something went wrong\n", err.Error())
}
