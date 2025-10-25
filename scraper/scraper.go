package scraper

import (
	"awesomeProject1/status"
	"context"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var Client *http.Client = &http.Client{
	Timeout: time.Second * 5,
	Transport: &http.Transport{
		MaxIdleConns:       100,
		MaxConnsPerHost:    100,
		IdleConnTimeout:    time.Second * 60,
		DisableCompression: false,
	},
}

func Run(urls []string, onResult func()) []Result {
	results := make([]Result, 0, len(urls))
	ch := make(chan Result)
	wg := sync.WaitGroup{}

	sem := make(chan struct{}, CountGoroutines)

	for _, url := range urls {
		status.SetStatus(url, "pending")
	}

	limiter := rate.NewLimiter(5, 2)

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			limiter.Wait(context.Background())

			sem <- struct{}{}
			defer func() { <-sem }()

			status.SetStatus(u, "in_progress")

			t := time.Now()
			status.SetTimeToAction(u, t)

			r := fetchUrlWithRetry(u, 3)
			status.SetStatus(u, "done")
			duration := time.Since(t)
			status.SetDuration(u, duration)

			if onResult != nil {
				onResult()
			}

			ch <- r
		}(url)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for r := range ch {
		results = append(results, r)
	}

	return results
}

func fetchUrl(url string) Result {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Add("User-Agent", "Scrapper/1.0")

	resp, err := Client.Do(req)

	if err != nil {
		return Result{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{URL: url, Error: resp.Status}
	}

	re := regexp.MustCompile(`<title>(.*?)</title>`)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}
	}
	matches := re.FindSubmatch(data)

	title := "N/A"
	if len(matches) > 1 {
		title = string(matches[1])
	}

	return Result{URL: url, Title: title, Error: resp.Status}
}

func fetchUrlWithRetry(url string, retryTimes int) Result {
	var result Result
	for i := 0; i <= retryTimes; i++ {
		result = fetchUrl(url)
		if result.Error == "" {
			return result
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return Result{URL: url, Error: result.Error}
}
