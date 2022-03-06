package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/mariaefi29/myhttp/md5_calculator"
)

type args struct {
	Workers int      `arg:"--parallel" help:"number of parallel workers" default:"10"`
	URLs    []string `arg:"positional"`
}

func main() {
	var args args
	arg.MustParse(&args)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	calculator := md5_calculator.Calculator{Client: &http.Client{Timeout: time.Minute}}
	results := calculator.CalcHashes(ctx, args.Workers, args.URLs)

	for _, result := range results {
		output := result.MD5Hash
		if result.Err != nil {
			output = fmt.Errorf("err: %s", result.Err).Error()
		}
		fmt.Println(result.InputURL, output)
	}
}
