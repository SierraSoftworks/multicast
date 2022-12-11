package multicast_test

import (
	"fmt"
	"sync"

	"github.com/SierraSoftworks/multicast/v2"
)

func ExampleReadme() {
	c := multicast.New[any]()
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		l := c.Listen()
		wg.Done()
		defer wg.Done()
		for msg := range l.C {
			fmt.Printf("Listener 1: %s\n", msg)
		}
		fmt.Println("Listener 1 Closed")
	}()

	go func() {
		l := c.Listen()
		wg.Done()
		defer wg.Done()
		for msg := range l.C {
			fmt.Printf("Listener 2: %s\n", msg)
		}
		fmt.Println("Listener 2 Closed")
	}()

	wg.Wait()
	wg.Add(2)
	c.C <- "Hello World!"
	c.Close()

	wg.Wait()

	// Unordered output:
	// Listener 1: Hello World!
	// Listener 2: Hello World!
	// Listener 1 Closed
	// Listener 2 Closed
}
