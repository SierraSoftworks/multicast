package multicast_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/SierraSoftworks/multicast/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleListener() {
	m := multicast.New[any]()
	l := m.Listen()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for msg := range l.C {
			fmt.Printf("Listener got: %#v\n", msg)
		}
		wg.Done()
	}()

	m.C <- "Hello!"
	m.Close()
	wg.Wait()

	// Output:
	// Listener got: "Hello!"
}

func ExampleNewListener() {
	s := make(chan interface{})

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		l := multicast.NewListener(s)
		fmt.Printf("Listener got: %s\n", <-l.C)
		wg.Done()
	}()

	s <- "Hello World!"
	close(s)
	wg.Wait()

	// Output:
	// Listener got: Hello World!
}

func ExampleListener_Chain() {
	s := make(chan interface{})

	wg := sync.WaitGroup{}
	wg.Add(1)

	l1 := multicast.NewListener(s)
	go func() {
		fmt.Printf("Listener 1: %s\n", <-l1.C)
		wg.Done()
	}()

	wg.Add(1)

	l2 := l1.Chain()
	go func() {
		fmt.Printf("Listener 2: %s\n", <-l2.C)
		wg.Done()
	}()

	s <- "Hello World!"
	close(s)

	wg.Wait()

	// Unordered Output:
	// Listener 1: Hello World!
	// Listener 2: Hello World!
}

func TestListener(t *testing.T) {
	Convey("Listener", t, func() {
		s := make(chan interface{}, 0)

		Convey("Single Listener", func() {
			l := multicast.NewListener(s)

			So(l, ShouldNotBeNil)
			So(l.C, ShouldNotBeNil)

			s <- "Hello"
			So(<-l.C, ShouldEqual, "Hello")

			Convey("Multiple Writes", func() {
				s <- "World"
				So(<-l.C, ShouldEqual, "World")
			})

			Convey("Closing", func() {
				close(s)
				_, ok := <-l.C
				So(ok, ShouldBeFalse)
			})
		})

		Convey("Chained Listeners", func() {
			l1 := multicast.NewListener(s)
			So(l1, ShouldNotBeNil)
			So(l1.C, ShouldNotBeNil)

			l2 := l1.Chain()
			So(l2, ShouldNotBeNil)
			So(l2.C, ShouldNotBeNil)

			Convey("Ordered Reads", func() {
				s <- "Hello"
				So(<-l1.C, ShouldEqual, "Hello")
				So(<-l2.C, ShouldEqual, "Hello")
			})

			Convey("Out of Order Reads", func() {
				s <- "Hello"
				So(<-l2.C, ShouldEqual, "Hello")
				So(<-l1.C, ShouldEqual, "Hello")
			})

			Convey("Closing", func() {
				close(s)

				_, ok := <-l1.C
				So(ok, ShouldBeFalse)

				_, ok = <-l2.C
				So(ok, ShouldBeFalse)
			})
		})
	})
}
