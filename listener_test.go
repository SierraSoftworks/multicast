package multicast_test

import (
	"fmt"
	"testing"

	"github.com/SierraSoftworks/multicast"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleListener() {
	m := multicast.New()
	l := m.Listen()

	go func() {
		for msg := range l.C {
			fmt.Printf("Listener got: %#v\n", msg)
		}
	}()

	m.C <- "Hello!"
}

func TestListener(t *testing.T) {
	Convey("Listener", t, func() {
		s := make(chan interface{}, 1)

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
