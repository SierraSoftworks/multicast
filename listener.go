package multicast

// Listener represents a listener which will receive messages
// from a channel.
type Listener struct {
	C <-chan interface{}
	f chan interface{}
}

// NewListener creates a new listener which will forward messages
// it receives on its f channel before exposing them on its C
// channel.
// You will very rarely need to use this method directly in your
// applications, prefer using From instead.
func NewListener(source <-chan interface{}) *Listener {
	out := make(chan interface{}, 0)
	l := &Listener{
		C: out,
	}

	go func() {
		for v := range source {
			if l.f != nil {
				l.f <- v
			}
			out <- v
		}

		if l.f != nil {
			close(l.f)
		}
		close(out)
	}()

	return l
}

// Chain is a shortcut which updates an existing listener to forward
// to a new listener and then returns the new listener.
// You will generally not need to make use of this method in your
// applications.
func (l *Listener) Chain() *Listener {
	f := make(chan interface{})
	l.f = f
	return NewListener(f)
}
