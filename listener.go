package multicast

// Listener represents a listener which will receive messages
// from a channel.
type Listener[T any] struct {
	C <-chan T
	f chan T
}

// NewListener creates a new listener which will forward messages
// it receives on its f channel before exposing them on its C
// channel.
// You will very rarely need to use this method directly in your
// applications, prefer using From instead.
func NewListener[T any](source <-chan T) *Listener[T] {
	out := make(chan T, 0)
	l := &Listener[T]{
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
func (l *Listener[T]) Chain() *Listener[T] {
	f := make(chan T)
	l.f = f
	return NewListener[T](f)
}
