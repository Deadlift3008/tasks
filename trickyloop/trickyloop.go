package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var max int32
	wg := sync.WaitGroup{}
	wg.Add(1000)

	for i := 1000; i > 0; i-- {
		go func() {
			defer wg.Done()

			if i%2 != 0 {
				return
			}

			for {
				localMax := atomic.LoadInt32(&max)

				if i > int(localMax) {
					result := atomic.CompareAndSwapInt32(&max, localMax, int32(i))

					if result {
						return
					}
				} else {
					return
				}
			}
		}()
	}

	wg.Wait()

	fmt.Println("max - ", max)
}
