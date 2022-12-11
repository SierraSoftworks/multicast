package multicast_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/SierraSoftworks/multicast/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleNew() {
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
	}()

	go func() {
		l := c.Listen()
		wg.Done()
		defer wg.Done()
		for msg := range l.C {
			fmt.Printf("Listener 2: %s\n", msg)
		}
	}()

	wg.Wait()
	wg.Add(2)
	c.C <- "Hello World!"
	c.Close()
	wg.Wait()

	// Unordered Output:
	// Listener 1: Hello World!
	// Listener 2: Hello World!
}

func ExampleFrom() {
	source := make(chan interface{})
	c := multicast.From(source)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		l := c.Listen()
		wg.Done()
		defer wg.Done()
		for msg := range l.C {
			fmt.Printf("Listener 1: %s\n", msg)
		}
	}()

	go func() {
		l := c.Listen()
		wg.Done()
		defer wg.Done()
		for msg := range l.C {
			fmt.Printf("Listener 2: %s\n", msg)
		}
	}()

	wg.Wait()
	wg.Add(2)
	source <- "Hello World!"
	close(source)
	wg.Wait()

	// Unordered Output:
	// Listener 1: Hello World!
	// Listener 2: Hello World!
}

func ExampleChannel_Close() {
	c := multicast.New[any]()
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		l := c.Listen()
		for range l.C {
		}
		fmt.Println("Listener closed")
		wg.Done()
	}()

	c.Close()
	wg.Wait()

	// Output:
	// Listener closed
}

func TestChannel(t *testing.T) {
	Convey("Channel", t, func() {
		Convey("Constructor", func() {
			c := multicast.New[any]()
			So(c, ShouldNotBeNil)
			So(c.C, ShouldNotBeNil)
		})

		Convey("From", func() {
			s := make(chan interface{})
			c := multicast.From(s)
			So(c, ShouldNotBeNil)
			So(c.C, ShouldNotBeNil)

			l := c.Listen()
			So(l, ShouldNotBeNil)

			go func() {
				s <- "Hello"
			}()
			So(<-l.C, ShouldEqual, "Hello")

			close(s)
			_, ok := <-l.C
			So(ok, ShouldBeFalse)
		})

		Convey("Listen", func() {
			c := multicast.New[any]()
			So(c, ShouldNotBeNil)

			l := c.Listen()
			So(l, ShouldNotBeNil)
			So(l.C, ShouldNotBeNil)

			go func() {
				c.C <- "Hello"
			}()
			So(<-l.C, ShouldEqual, "Hello")

			l2 := c.Listen()
			So(l2, ShouldNotBeNil)
			So(l2.C, ShouldNotBeNil)

			go func() {
				c.C <- "World"
			}()
			So(<-l.C, ShouldEqual, "World")
			So(<-l2.C, ShouldEqual, "World")
		})

		Convey("Close", func() {
			c := multicast.New[any]()
			So(c, ShouldNotBeNil)

			l := c.Listen()
			So(l, ShouldNotBeNil)

			c.Close()

			_, ok := <-l.C
			So(ok, ShouldBeFalse)
		})
	})
}
