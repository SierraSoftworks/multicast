package multicast

import "sync"

// Channel represents a multicast channel container which provides
// a writable channel C and allows multiple listeners to be connected.
type Channel struct {
	// C is a writable channel on which you can send messages which will
	// be delivered to all connected listeners.
	C chan<- interface{}

	c chan interface{}
	l *Listener
	m sync.Mutex
}

// New creates a new multicast channel container which can have listeners
// connected to it and messages sent via its C channel property.
func New() *Channel {
	c := make(chan interface{})

	return From(c)
}

// From creates a new multicast channel which exposes messages it receives
// on the provided channel to all connected listeners.
func From(c chan interface{}) *Channel {
	return &Channel{
		C: c,
		c: c,
	}
}

// Listen returns a new listener instance attached to this channel.
// Each listener will receive a single instance of each message sent
// to the channel.
func (c *Channel) Listen() *Listener {
	c.m.Lock()
	defer c.m.Unlock()

	if c.l == nil {
		c.l = NewListener(c.c)
	} else {
		c.l = c.l.Chain()
	}

	return c.l
}

// Close is a convinience function for closing the top level channel.
// You may also close the channel directly by using `close(c.C)`.
func (c *Channel) Close() {
	c.m.Lock()
	defer c.m.Unlock()

	close(c.c)
}
