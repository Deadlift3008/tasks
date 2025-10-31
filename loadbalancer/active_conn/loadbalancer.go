package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Request(url string) {
	time.Sleep(time.Duration(rand.Int31n(10)) * time.Second)
	fmt.Printf("Request done for %v\n", url)
}

type Node struct {
	url         string
	active_conn int
}

type LoadBalancer struct {
	nodes []Node
	mu    sync.RWMutex
}

func (l *LoadBalancer) findMinConnNode() *Node {
	minNode := &l.nodes[0]

	for i := 0; i < len(l.nodes); i++ {
		if l.nodes[i].active_conn < minNode.active_conn {
			minNode = &l.nodes[i]
		}
	}

	return minNode
}

func (l *LoadBalancer) DoRequest() {
	l.mu.Lock()
	node := l.findMinConnNode()
	node.active_conn++
	l.mu.Unlock()

	Request(node.url)

	l.mu.Lock()
	node.active_conn--
	l.mu.Unlock()
}

func (l *LoadBalancer) GetStats() map[string]int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	stats := make(map[string]int)
	for _, node := range l.nodes {
		stats[node.url] = node.active_conn
	}
	return stats
}

func NewLoadBalancer(urls []string) *LoadBalancer {
	nodes := make([]Node, len(urls))

	for index, url := range urls {
		nodes[index] = Node{url: url}
	}

	return &LoadBalancer{
		nodes: nodes,
	}
}

func main() {
	lb := NewLoadBalancer([]string{
		"http://testurl.ru",
		"https://google.ru",
		"https://yandex.ru",
		"https://youtube.ru",
	})

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lb.DoRequest()

			if i%3 == 0 {
				fmt.Println(lb.GetStats())
			}
		}()
	}

	wg.Wait()
}
