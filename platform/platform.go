package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.SetLimit(2)

	for i := 0; i < 5; i++ {
		i := i
		g.Go(func() error {
			if i == 3 {
				return errors.New("EBANA VROT")
			}

			if ctx.Err() != nil {
				return nil
			}

			fmt.Printf("Goroutine %v\n", i)
			return nil
		})
	}

	err := g.Wait()

	if err != nil {
		fmt.Println("OPANA")
		fmt.Println(err.Error())
	}
}
