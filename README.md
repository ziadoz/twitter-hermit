# Twitter Hermit
Become a social media hermit. A Twitter utility that deletes old tweets, retweets and favourites.

# Usage
You must ensure your [Twitter OAuth credentials](https://developer.twitter.com/apps) are available as environment variables: `TWITTER_CONSUMER_KEY`, `TWITTER_CONSUMER_SECRET`, `TWITTER_ACCESS_TOKEN` and `TWITTER_ACCESS_TOKEN_SECRET`

Download the latest release and ensure the binary is available on your `$PATH` and then run it with the required `--max-age` argument:
```
twitter-hermit --max-age="-1month"
```

The `--max-age=[duration]` argument determines how old tweets have to be before they'll be deleted. This string is parsed using the [tparse](https://github.com/karrick/tparse) package, and some examples include: `-1day`, `-2weeks`, `-3months` and `-4years`. You can find more details on the format here: https://github.com/karrick/tparse#addduration
```
twitter-hermit --max-age="-1month"
```

You can also perform a dry run with the `--dry-run` flag which will only show a summary of the tweets, retweets and favourites it would delete:
```
twitter-hermit --max-age="-1month" --dry-run
```

If you don't care about the summary you can use the `--silent` flag to suppress any output (this does not work in dry run mode):
```
twitter-hermit --max-age="-1month" --silent
```

Hermit can also save JSON, media and links from any tweets into a file before it deletes them with the `--save-dir=[dir]` argument, and the `--save-json`, `--save-media` and `--save-links` flags:
```
twitter-hermit --max-age="-1month" --save-dir=./tweets --save-json --save-media --save-links
```

For convience a `env.sh.sample` file is provided. You should rename this to `env.sh`, enter your credentials into it, specify your preferred flags, and then run it: `./env.sh`

## Notes
This project was loosely inspired by several existing projects. I highly recommend checking them out, they're amazing:

- https://mkaz.blog/misc/automate-deleting-your-tweets-with-a-raspberry-pi/
- https://vickylai.com/verbose/delete-old-tweets-ephemeral/
- https://github.com/adamdrake/harold
