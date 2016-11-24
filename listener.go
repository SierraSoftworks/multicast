package multicast

// Listener represents a listener which will receive messages
// from a channel.
type Listener struct {
	C <-chan interface{}
	f chan interface{}
}

func newListener(source <-chan interface{}) *Listener {
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

func (l *Listener) chain() *Listener {
	f := make(chan interface{})
	l.f = f
	return newListener(f)
}
