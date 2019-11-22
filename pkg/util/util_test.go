package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestStripNewlines(t *testing.T) {
	is := is.New(t)
	have := "foo\nbar\r\nbaz\tqux\n\n"
	want := "foo bar baz\tqux"
	is.Equal(StripNewlines(have), want)
}
