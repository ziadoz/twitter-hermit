# Twitter Hermit
Become a social media hermit. A Twitter utility that deletes old tweets, retweets and favourites.

# Usage
You must ensure your [Twitter OAuth credentials](https://developer.twitter.com/apps) are available as environment variables: `TWITTER_CONSUMER_KEY`, `TWITTER_CONSUMER_SECRET`, `TWITTER_ACCESS_TOKEN` and `TWITTER_ACCESS_TOKEN_SECRET`

Download the latest release and ensure the binary is available on your `$PATH` and then run it with the required `--max-age` flag:
```
twitter-hermit --max-age="1 month"
```

The `--max-age` flag supports the following suffixes: `day`, `week`, `month` and `year`. You can also use the pluralised versions (e.g. `2 days`, `3 weeks`, `4 months`, `5 years`).

You can also perform a dry run with the `--dry-run` flag which will only show a summary of the tweets, retweets and favourites it would delete:
```
twitter-hermit --max-age="1 month" --dry-run
```

If you don't care about the summary you can use the `--silent` flag to suppress any output (this does not work in dry run mode):
```
twitter-hermit --max-age="1 month" --silent
```

Hermit can also extract links from any tweets into a file before it deletes them with the `--extract-links` flags:
```
twitter-hermit --max-age="1 month" --extract-links="./links.txt"
```

For convience a `env.sh.sample` file is provided. You should rename this to `env.sh`, enter your credentials into it, specify your preferred flags, and then run it: `./env.sh`

## Notes
This project was loosely inspired by several existing projects. I highly recommend checking them out, they're amazing:

- https://mkaz.blog/misc/automate-deleting-your-tweets-with-a-raspberry-pi/
- https://vickylai.com/verbose/delete-old-tweets-ephemeral/
- https://github.com/adamdrake/harold
