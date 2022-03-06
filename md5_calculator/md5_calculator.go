package md5_calculator

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Calculator struct {
	Client *http.Client
}

type CalcHashResult struct {
	InputURL string
	MD5Hash  string
	Err      error
}

func (c Calculator) CalcHashes(ctx context.Context, workers int, urls []string) []CalcHashResult {
	urlCh := make(chan string)

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wg.Add(workers)

	go func() {
		defer close(urlCh)
		writeURL(ctx, urlCh, urls)
	}()

	results := make([]CalcHashResult, 0, len(urls))
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			res := c.processURLs(ctx, urlCh)
			mu.Lock()
			results = append(results, res...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return results
}

func (c Calculator) CalcHash(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %s", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read resp body: %s", err)
	}

	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
}

func (c Calculator) processURLs(ctx context.Context, urlCh chan string) []CalcHashResult {
	results := make([]CalcHashResult, 0)
	for url := range urlCh {
		hash, err := c.CalcHash(ctx, url)
		results = append(results, CalcHashResult{
			InputURL: url,
			MD5Hash:  hash,
			Err:      err,
		})
	}

	return results
}

func writeURL(ctx context.Context, urlCh chan string, input []string) {
	for _, url := range input {
		select {
		case <-ctx.Done():
			return
		case urlCh <- url:
		}
	}
	return
}
