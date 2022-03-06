package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/mariaefi29/myhttp/md5_calculator"
)

type args struct {
	Workers int      `arg:"--parallel" help:"number of parallel workers" default:"10"`
	Input   []string `arg:"positional"`
}

func main() {
	var args args
	arg.MustParse(&args)

	urlCh := make(chan string)
	ctx := context.Background()

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wg.Add(args.Workers)

	go func() {
		defer close(urlCh)
		writeURL(ctx, urlCh, args.Input)
	}()

	calculator := md5_calculator.Calculator{
		Client: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}
	results := make([]result, 0, len(args.Input))
	for i := 0; i < args.Workers; i++ {
		go func() {
			defer wg.Done()
			res := processURLs(ctx, calculator, urlCh)
			mu.Lock()
			results = append(results, res...)
			mu.Unlock()
		}()
	}

	signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	wg.Wait()

	for _, result := range results {
		output := result.MD5Hash
		if result.err != nil {
			output = fmt.Errorf("err: %s", result.err).Error()
		}
		fmt.Println(result.inputURL, output)
	}
}

func writeURL(ctx context.Context, urlCh chan string, input []string) {
	for _, url := range input {
		const prefix = "https://"
		if !strings.HasPrefix(url, prefix) {
			url = prefix + url
		}
		select {
		case <-ctx.Done():
			return
		case urlCh <- url:
		}
	}

	return
}

type result struct {
	inputURL string
	MD5Hash  string
	err      error
}

func processURLs(ctx context.Context, calculator md5_calculator.Calculator, urlCh chan string) []result {
	results := make([]result, 0)
	for url := range urlCh {
		hash, err := calculator.CalcMD5Hash(ctx, url)
		results = append(results, result{
			inputURL: url,
			MD5Hash:  hash,
			err:      err,
		})
	}

	return results
}
