package util_test

import (
	"testing"

	"github.com/ziadoz/twitter-hermit/pkg/util"

	"github.com/matryer/is"
)

func TestStripNewlines(t *testing.T) {
	is := is.New(t)

	have := "foo\nbar\r\nbaz\tqux\n\n"
	want := "foo bar baz\tqux"

	is.Equal(util.StripNewlines(have), want)
}
