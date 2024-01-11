package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	done := make(chan struct{})
	single := make(chan interface{})
	for i := range channels {
		go func(ch <-chan interface{}) {
			select {
			case v := <-ch:
				close(done)
				single <- v
			case <-done:
				return
			}

		}(channels[i])
	}
	<-done
	return single
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Second),
		sig(5*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
	)

	fmt.Printf("fone after %v", time.Since(start))

}
