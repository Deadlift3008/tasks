package ratelimiter

import (
	"time"
)

type Request struct {
	Time time.Time
}

type RateLimiter struct {
	rpsLimit int
	requests []Request
}

func (r *RateLimiter) IsAllow() bool {
	now := time.Now()

	secondAgo := now.Add(time.Duration(-1) * time.Second)

	lastRequests := make([]Request, 0, r.rpsLimit)

	for _, request := range r.requests {
		if request.Time.After(secondAgo) {
			lastRequests = append(lastRequests, Request{request.Time})
		}
	}

	doneRequests := len(lastRequests)

	lastRequests = append(lastRequests, Request{now})

	r.requests = lastRequests

	return doneRequests < r.rpsLimit
}

func NewRateLimiter(rpsLimit int) RateLimiter {
	return RateLimiter{
		rpsLimit: rpsLimit,
	}
}
