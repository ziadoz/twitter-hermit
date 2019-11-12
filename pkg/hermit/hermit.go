package hermit

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ziadoz/twitter-hermit/pkg/saver"
	"github.com/ziadoz/twitter-hermit/pkg/twitter"
	"github.com/ziadoz/twitter-hermit/pkg/util"
)

const batchSize = 200

type Destroyer struct {
	MaxAge     time.Time         // The max age to filter out tweets for deletion.
	DryRun     bool              // Whether or not the deletion should be a dry run.
	Output     io.Writer         // Output is written to this.
	TweetSaver *saver.TweetSaver // Handle saving tweet content.
}

func (d *Destroyer) Destroy(repo twitter.Repository) error {
	header := fmt.Sprintf("Destroying %s", strings.Title(repo.Description()))
	fmt.Fprintln(d.Output, header)
	fmt.Fprintln(d.Output, strings.Repeat("=", utf8.RuneCountInString(header)))

	var maxID int64
	var deletedCount int

	for {
		tweets, err := repo.Get(twitter.QueryParams{Count: batchSize, MaxID: maxID})
		if err != nil {
			return fmt.Errorf("could not get user %s: %s\n", repo.Description(), err)
		}

		if len(tweets) == 0 {
			break // We're done deleting.
		}

		filteredTweets := twitter.FilterTweets(tweets, d.MaxAge)
		if len(filteredTweets) == 0 {
			maxID = twitter.GetMaxID(tweets) - 1
			continue
		}

		if d.TweetSaver != nil {
			for _, tweet := range filteredTweets {
				if err := d.TweetSaver.Save(tweet); err != nil {
					return fmt.Errorf("could not save tweet '%d' content: %s", tweet.ID, err)
				}
			}
		}

		if !d.DryRun {
			err = repo.Destroy(filteredTweets)
			if err != nil {
				return fmt.Errorf("could not get user %s: %s\n", repo.Description(), err)
			}
		}

		for _, tweet := range filteredTweets {
			createdAt, _ := tweet.CreatedAtTime()
			fmt.Fprintf(
				d.Output,
				" - ID:   %d\n   Date: %s\n   Text: %s\n",
				tweet.ID,
				createdAt.Format("2 Jan 2006 03:04pm"),
				util.StripNewlines(tweet.Text),
			)
		}

		deletedCount += len(filteredTweets)
		maxID = twitter.GetMaxID(tweets) - 1
	}

	if deletedCount > 0 {
		fmt.Fprintf(d.Output, "Deleted %d %s successfully!\n", deletedCount, repo.Description())
	} else {
		fmt.Fprintf(d.Output, "No %s needed deleting.\n", repo.Description())
	}

	return nil
}
