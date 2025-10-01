package ratelimiter

import (
	"fmt"
	"testing"
	"time"
)

func TestRatelimiter(t *testing.T) {
	rl := NewRateLimiter(1)

	first := rl.IsAllow()
	time.Sleep(time.Duration(1) * time.Second)
	second := rl.IsAllow()
	time.Sleep(time.Duration(1) * time.Second)
	third := rl.IsAllow()
	time.Sleep(time.Duration(1) * time.Second)
	fourth := rl.IsAllow()

	fmt.Println(first)
	fmt.Println(second)
	fmt.Println(third)
	fmt.Println(fourth)
}
