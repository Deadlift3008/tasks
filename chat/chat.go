package main

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

type IChat interface {
	Send(message string) int
	Receive(id int) string
}

type Shard struct {
	messages map[int]string
	mu       sync.RWMutex
	lastId   int
}

type Chat struct {
	shards         map[int]*Shard
	messageCounter int32
}

func (c *Chat) Send(message string) int {
	mCount := atomic.AddInt32(&c.messageCounter, 1)

	shardNumber := int(mCount) % len(c.shards)

	shard := c.shards[shardNumber]

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.lastId += len(c.shards)
	shard.messages[shard.lastId] = message

	return shard.lastId
}

func (c *Chat) Receive(id int) string {
	shardNumber := id % len(c.shards)

	shard := c.shards[shardNumber]

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	message := shard.messages[id]

	return message
}

func NewChat(shardsCount int) *Chat {
	shards := make(map[int]*Shard)

	for i := 0; i < shardsCount; i++ {
		shards[i] = &Shard{
			messages: make(map[int]string),
			lastId:   i,
		}
	}

	return &Chat{
		shards: shards,
	}
}

func main() {
	chat := NewChat(5)

	ids := []int{}
	wg := sync.WaitGroup{}
	var mu sync.Mutex

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := chat.Send("message num: " + strconv.Itoa(i))
			mu.Lock()
			ids = append(ids, id)
			mu.Unlock()
		}()
	}

	wg.Wait()

	for _, msgId := range ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(chat.Receive(msgId))
		}()
	}

	wg.Wait()
}
