package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type User struct {
	Name string
}

func main() {
	t := time.Now()
	names, err := Do(context.Background(), []User{
		{"Paul"}, {"Jack"}, {"Jack"}, {"Mike"},
	})
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Result: ", names)
	fmt.Println("Elapsed time: ", time.Since(t))
}

// fetch что-то дeлaeм по сeти.
func fetch(_ context.Context, u User) (string, error) {
	time.Sleep(time.Millisecond * 10) // имитaция зaдeржки
	return u.Name, nil
}

func Do(ctx context.Context, users []User) (map[string]int64, error) {
	names := make(map[string]int64, len(users))

	wg := sync.WaitGroup{}

	resultChan := make(chan string)
	workChan := make(chan User)

	workersNum := runtime.NumCPU() * 3

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for _, u := range users {
			workChan <- u
		}
		close(workChan)
	}()

	wg.Add(workersNum)
	for range workersNum {
		go func() {
			defer wg.Done()

			for u := range workChan {
				if ctx.Err() != nil {
					return
				}

				res, err := fetch(ctx, u)

				if err != nil {
					fmt.Println(err)
					cancel()
				}

				resultChan <- res
			}
		}()
	}

	go func() {
		for result := range resultChan {
			names[result]++
		}
	}()

	wg.Wait()
	close(resultChan)

	return names, nil
}
