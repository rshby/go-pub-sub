package test

import (
	"fmt"
	"sync"
	"testing"
)

func TestDebugGoroutine(t *testing.T) {
	fmt.Println("test")

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		fmt.Println("di dalam goroutine")
	}(wg)

	wg.Wait()
	fmt.Println("selesai")
}
