package main

import (
	"fmt"
	"sync"
	"time"
)

func checkUrl1(urls []string) {
	startTime := time.Now()
	for _, url := range urls {
		// URLにアクセスしてステータスコードを返す
		result := check(url, time.Second*2)
		if result.Err != nil {
			fmt.Println("error")
			return
		}

		fmt.Println(result.StatusCode, result.URL)
	}
	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))

}

func checkUrl2(urls []string) {
	startTime := time.Now()
	ch := make(chan Result, len(urls))
	for _, url := range urls {
		// URLにアクセスしてステータスコードを返す
		go func(u string) {
			result := check(u, time.Second*2)
			ch <- result // Result 構造体をそのまま送る
		}(url)
	}

	for range urls {
		r := <-ch
		if r.Err != nil {
			fmt.Printf("error: %s: %v\n", r.URL, r.Err)
		} else {
			fmt.Printf("%d %s\n", r.StatusCode, r.URL)
		}
	}

	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))

}
func checkUrl3(urls []string) {
	startTime := time.Now()
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := check(u, time.Second*2)
			if result.Err != nil {
				fmt.Println("error", result.Err)
			} else {
				fmt.Printf("%d %s\n", result.StatusCode, result.URL)
			}
		}(url)
	}

	wg.Wait()

	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))
}

func checkUrl4(urls []string, limit int) {
	startTime := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			// URLにアクセスしてステータスコードを返す
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			result := check(u, time.Second*2)
			if result.Err != nil {
				fmt.Printf("error: %s: %v\n", result.URL, result.Err)
			} else {
				fmt.Printf("%d: %s\n", result.StatusCode, result.URL)
			}
		}(url)
	}

	wg.Wait()

	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))

}
func main() {

	urls := []string{"https://google.com", "https://golang.org", "https://google.com", "https://golang.org", "https://google.com", "https://golang.org", "https://google.com", "https://golang.org", "https://google.com", "https://golang.org", "https://google.com", "https://golang.org"}

	// 直列で実行した際の実行時間を測定
	fmt.Println("checkUrl1 start")
	checkUrl1(urls)

	// 並列(channel)で実行した際の実行時間を測定
	fmt.Println("checkUrl2 start")
	checkUrl2(urls)

	// 並列(waitGroup)で実行した際の実行時間を測定
	fmt.Println("checkUrl3 start")
	checkUrl3(urls)

	// 並列(channel 同時実行数の制御あり)で実行した際の実行時間を測定
	fmt.Println("checkUrl4 start")
	checkUrl4(urls, 3)
}
