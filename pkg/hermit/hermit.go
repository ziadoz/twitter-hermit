package hermit

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ziadoz/twitter-hermit/pkg/saver"
	"github.com/ziadoz/twitter-hermit/pkg/twitter"
	"github.com/ziadoz/twitter-hermit/pkg/util"

	"github.com/olekukonko/tablewriter"
)

const batchSize = 200

type Destroyer struct {
	MaxAge     time.Time         // The max age to filter out tweets for deletion.
	DryRun     bool              // Whether or not the deletion should be a dry run.
	Output     io.Writer         // Output is written to this.
	TweetSaver *saver.TweetSaver // Handle saving tweet content.
}

func (d *Destroyer) Destroy(repo twitter.Repository) error {
	table := tablewriter.NewWriter(d.Output)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.SetHeader([]string{"ID", "Date", "Tweet"})

	var maxID int64
	var deletedCount int

	for {
		tweets, err := repo.Get(twitter.QueryParams{Count: batchSize, MaxID: maxID})
		if err != nil {
			return fmt.Errorf("could not get user %s: %s\n", repo.Description(), err)
		}

		if len(tweets) == 0 {
			break
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

		rows := make([][]string, 0, len(filteredTweets))
		for _, tweet := range filteredTweets {
			createdAt, _ := tweet.CreatedAtTime()
			rows = append(rows, []string{
				strconv.FormatInt(tweet.ID, 10),
				createdAt.Format("2 Jan 2006 03:04pm"),
				util.StripNewlines(tweet.Text),
			})
		}

		table.AppendBulk(rows)

		deletedCount += len(filteredTweets)
		maxID = twitter.GetMaxID(tweets) - 1
	}

	if deletedCount > 0 {
		table.SetFooter([]string{"", repo.Description(), strconv.Itoa(deletedCount)})
		table.Render()
	}

	return nil
}
