package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func fetch(url string) error {
	fmt.Printf("fetch %v", url)
	time.Sleep(time.Second * 2)
	return nil
}

func MassFetch(urls []string, errorLimit int32) error {
	workersNum := runtime.NumCPU() * 3
	urlChan := make(chan string)
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	var errorsCount int32

	wg.Add(workersNum)
	for range workersNum {
		go func() {
			defer wg.Done()

			for {
				if ctx.Err() != nil {
					return
				}

				select {
				case url := <-urlChan:
					err := fetch(url)

					if err != nil {
						fmt.Printf("url %v - no OK", url)

						currentCountErr := atomic.AddInt32(&errorsCount, 1)

						if currentCountErr > errorLimit {
							cancel(errors.New("exceeded error limit"))
						}
					} else {
						fmt.Printf("url %v - 200 OK", url)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		defer close(urlChan)

		for _, url := range urls {
			select {
			case urlChan <- url:
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()

	return ctx.Err()
}

func main() {
	var urls = []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://eeee",
	}

	// делает запросы по урлам, при получении 2 ошибок - стопится
	MassFetch(urls, 2)
}
