package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

func Request(url string) {
	time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
	fmt.Printf("Request done for %v", url)
}

type LoadBalancer struct {
	urls    []string
	counter int32
	lenUrls int
}

func (l *LoadBalancer) DoRequest() {
	nextCounter := atomic.AddInt32(&l.counter, 1)

	nextUrlIndex := nextCounter % int32(l.lenUrls)

	Request(l.urls[nextUrlIndex])
}

func NewLoadBalancer(urls []string) *LoadBalancer {
	return &LoadBalancer{
		urls:    urls,
		lenUrls: len(urls),
	}
}

func main() {

}
