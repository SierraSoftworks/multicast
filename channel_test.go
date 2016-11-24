package multicast_test

import (
	"fmt"
	"testing"

	"github.com/SierraSoftworks/multicast"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleChannel() {
	c := multicast.New()

	go func() {
		l := c.Listen()
		for msg := range l.C {
			fmt.Printf("Listener 1: %s\n", msg)
		}
	}()

	go func() {
		l := c.Listen()
		for msg := range l.C {
			fmt.Printf("Listener 2: %s\n", msg)
		}
	}()

	c.C <- "Hello World!"
	c.Close()
}

func ExampleFromChannel() {
	source := make(chan interface{})
	c := multicast.From(source)

	go func() {
		l := c.Listen()
		for msg := range l.C {
			fmt.Printf("Listener 1: %s\n", msg)
		}
	}()

	go func() {
		l := c.Listen()
		for msg := range l.C {
			fmt.Printf("Listener 2: %s\n", msg)
		}
	}()

	source <- "Hello World!"
	close(source)
}

func TestChannel(t *testing.T) {
	Convey("Channel", t, func() {
		Convey("Constructor", func() {
			c := multicast.New()
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
			c := multicast.New()
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
			c := multicast.New()
			So(c, ShouldNotBeNil)

			l := c.Listen()
			So(l, ShouldNotBeNil)

			c.Close()

			_, ok := <-l.C
			So(ok, ShouldBeFalse)
		})
	})
}
