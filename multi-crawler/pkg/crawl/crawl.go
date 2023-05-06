package pkg

import (
	"fmt"
	"sync"
	"sync/atomic"

	"extractHTTP"
)

// Crawler is a struct that contains the crawler's data
type Crawler struct {
	// URL to be crawled
	StartURL string
	// MaxDepth the maximum number of links to be crawled
	MaxDepth int
	// Maxconcurrency the maximum number of concurrent requests
	MaxConcurrency int
}

// Crawl the URL
func Crawl(config Crawler) ([]interface{}, error) {
	var (
		visited 			sync.Map
		queue 			=	make(chan string)
		results  		= make(chan []interface{}, config.MaxConcurrency)
		waitGroup 		sync.WaitGroup
	)

	visited.Store(config.StartURL, struct{}{})
	queue <- config.StartURL
	activeWokers := int32(0)

	processURL := func(url string, depth int) {
		atomic.AddInt32(&activeWokers, 1)
		defer atomic.AddInt32(&activeWokers, -1)
		defer waitGroup.Done()

		content, err := extractHTTP.FetchURLContent(url)
		if err != nil {
			return
		}
		fmt.Println(content)

		links, info := extractHTTP.ExtractLinksAndInfo(content)
		results <- info
		fmt.Println(links)
		fmt.Println(info)

		if depth < config.MaxDepth {
			for _, link := range links {
				if _, ok := visited.LoadOrStore(link, struct{}{}); !ok {
					waitGroup.Add(1)
					queue <- link
				}
			}
		}
	}

	for i := 0; i < config.MaxConcurrency; i++ {
		go func() {
			for {
				select {
				case url, ok := <-queue:
					if !ok {
						return
					}
					depth := 0 // Replace with actual depth calculation logic
					processURL(url, depth)
				default:
					if atomic.LoadInt32(&activeWokers) == 0 {
						return
					}
				}
			}
		}()
	}

	waitGroup.Wait()
	close(queue)
	close(results)

	var finalResults []interface{}
	for result := range results {
		finalResults = append(finalResults, result...)
	}

	return finalResults, nil
}