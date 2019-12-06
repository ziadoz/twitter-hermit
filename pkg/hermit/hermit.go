package hermit

import (
	"fmt"
	"io"
	"time"

	"github.com/ziadoz/twitter-hermit/pkg/data"
	"github.com/ziadoz/twitter-hermit/pkg/saver"

	"github.com/dustin/go-humanize"
)

const DefaultBatchSize = 200

type Destroyer struct {
	BatchSize  int               // The number of tweets to fetch per query.
	MaxAge     time.Time         // The max age to filter out tweets for deletion.
	DryRun     bool              // Whether or not the deletion should be a dry run.
	Output     io.Writer         // Output is written to this.
	TweetSaver *saver.TweetSaver // Handle saving tweet content.
}

func (d *Destroyer) Destroy(repo data.Repository) error {
	var maxID int64
	var deletedCount int

	for {
		tweets, err := repo.Get(data.QueryParams{Count: d.BatchSize, MaxID: maxID})
		if err != nil {
			return fmt.Errorf("could not get %s: %s\n", repo.Description(), err)
		}

		if len(tweets) == 0 {
			break
		}

		filteredTweets := data.FilterTweets(tweets, d.MaxAge)
		if len(filteredTweets) == 0 {
			maxID = data.GetMaxID(tweets) - 1
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
				return fmt.Errorf("could not destroy %s: %s\n", repo.Description(), err)
			}
		}

		deletedCount += len(filteredTweets)
		maxID = data.GetMaxID(tweets) - 1
	}

	if deletedCount > 0 {
		fmt.Fprintf(d.Output, "• Deleted %s %s.\n", humanize.Comma(int64(deletedCount)), repo.Description())
	}

	return nil
}
