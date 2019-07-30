package links

import (
	"net/http"
	"strings"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"mvdan.cc/xurls/v2"
)

// Extract links from tweet text.
func Extract(tweets []twitter.Tweet) []string {
	links := []string{}
	for _, tweet := range tweets {
		for _, link := range xurls.Strict().FindAllString(tweet.Text, -1) {
			links = append(links, link)
		}
	}
	return links
}

// Follow Twitter short URLs to get the actual URL.
func FollowRedirects(links []string) []string {
	redirected := make(chan string)
	sem := make(chan struct{}, 20)
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(redirected)
	}()

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			sem <- struct{}{}
			resp, err := http.Get(link)
			<-sem

			if err != nil {
				return
			}

			url := resp.Request.URL.String()
			if resp.StatusCode == http.StatusOK && !strings.HasPrefix(url, "https://t.co/") {
				redirected <- resp.Request.URL.String()
			}
		}(link)
	}

	results := []string{}
	for link := range redirected {
		results = append(results, link)
	}
	return results
}
