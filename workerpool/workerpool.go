package main

import (
	"fmt"
	"sync"
)

type Action = func()

type Workerpool struct {
	actions []Action
	wg      sync.WaitGroup
	limit   int
}

func (w *Workerpool) Start() {
	w.wg.Add(w.limit)

	actionChan := make(chan Action)

	for i := 0; i < w.limit; i++ {
		go func() {
			defer w.wg.Done()

			for action := range actionChan {
				action()
			}
		}()
	}

	go func() {
		for i := 0; i < len(w.actions); i++ {
			actionChan <- w.actions[i]
		}
		close(actionChan)
	}()

	w.wg.Wait()
}

func NewWorkerpool(actions []Action, limit int) *Workerpool {
	return &Workerpool{
		actions: actions,
		limit:   limit,
	}
}

func main() {
	actions := []Action{func() {
		fmt.Println("first")
	}, func() {
		fmt.Println("second")
	}, func() {
		fmt.Println("third")
	}, func() {
		fmt.Println("fourth")
	}, func() {
		fmt.Println("fifth")
	}}

	wp := NewWorkerpool(actions, 2)
	wp.Start()
}
