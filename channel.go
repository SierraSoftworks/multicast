package multicast

import "sync"

// Channel represents a multicast channel container which provides
// a writable channel C and allows multiple listeners to be connected.
type Channel[T any] struct {
	// C is a writable channel on which you can send messages which will
	// be delivered to all connected listeners.
	C chan<- T

	c chan T
	l *Listener[T]
	m sync.Mutex
}

// New creates a new multicast channel container which can have listeners
// connected to it and messages sent via its C channel property.
func New[T any]() *Channel[T] {
	c := make(chan T)

	return From(c)
}

// From creates a new multicast channel which exposes messages it receives
// on the provided channel to all connected listeners.
func From[T any](c chan T) *Channel[T] {
	return &Channel[T]{
		C: c,
		c: c,
	}
}

// Listen returns a new listener instance attached to this channel.
// Each listener will receive a single instance of each message sent
// to the channel.
func (c *Channel[T]) Listen() *Listener[T] {
	c.m.Lock()
	defer c.m.Unlock()

	if c.l == nil {
		c.l = NewListener[T](c.c)
	} else {
		c.l = c.l.Chain()
	}

	return c.l
}

// Close is a convenience function for closing the top level channel.
// You may also close the channel directly by using `close(c.C)`.
func (c *Channel[T]) Close() {
	c.m.Lock()
	defer c.m.Unlock()

	close(c.c)
}
