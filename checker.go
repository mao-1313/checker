package main

import (
	"context"
	"net/http"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Err        error
}

func check(url string, timeout time.Duration) Result {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Err: err}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{URL: url, Err: err}
	}
	defer resp.Body.Close()
	return Result{URL: url, StatusCode: resp.StatusCode}
}
