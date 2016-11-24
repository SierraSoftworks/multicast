package multicast

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListener(t *testing.T) {
	Convey("Listener", t, func() {
		s := make(chan interface{}, 1)

		Convey("Single Listener", func() {
			l := newListener(s)

			So(l, ShouldNotBeNil)
			So(l.C, ShouldNotBeNil)
			So(l.f, ShouldBeNil)

			s <- "Hello"
			So(<-l.C, ShouldEqual, "Hello")

			Convey("Closing", func() {
				close(s)
				_, ok := <-l.C
				So(ok, ShouldBeFalse)
			})
		})

		Convey("Chained Listeners", func() {
			l1 := newListener(s)
			So(l1, ShouldNotBeNil)
			So(l1.C, ShouldNotBeNil)
			So(l1.f, ShouldBeNil)

			l2 := l1.chain()
			So(l2, ShouldNotBeNil)
			So(l2.C, ShouldNotBeNil)
			So(l2.f, ShouldBeNil)
			So(l1.f, ShouldNotBeNil)

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
