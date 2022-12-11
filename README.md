# Multicast [![Build Status](https://travis-ci.org/SierraSoftworks/multicast.svg?branch=master)](https://travis-ci.org/SierraSoftworks/multicast) [![GoDoc](https://godoc.org/github.com/SierraSoftworks/multicast?status.svg)](https://godoc.org/github.com/SierraSoftworks/multicast/v2)
**Multi-subscriber channels for Golang**

The multicast module provides single-writer, multiple-reader semantics around Go channels.
It attempts to maintain semantics similar to those offered by standard Go channels while
guaranteeing parallel delivery (slow consumers won't hold up delivery to other listeners)
and guaranteeing delivery to all registered listeners when a message is published.

## Features

 - **Simple API** if you know how to use channels then you'll be able to use multicast
   without learning anything special.
 - **Similar Semantics** to core channels mean that your approach to reasoning around
   channels doesn't need to change.
 - **Low Overheads** with linear memory growth as your number of listeners increases
   and no buffer overheads.
 - **Generics Support** when using the `v2` library, allowing you to statically validate
   the types of messages you're sending and receiving.

## Example

```go
import (
    "fmt"

    "github.com/SierraSoftworks/multicast/v2"
)

func main() {
    c := multicast.New[string]()

	go func() {
		l := c.Listen()
		for msg := range l.C {
			fmt.Printf("Listener 1: %s\n", msg)
		}
        fmt.Println("Listener 1 Closed")
	}()

	go func() {
		l := c.Listen()
		for msg := range l.C {
			fmt.Printf("Listener 2: %s\n", msg)
		}
        fmt.Println("Listener 2 Closed")
	}()

	c.C <- "Hello World!"
	c.Close()
}
```

## Architecture
Multicast implements its channels as a linked list of listeners which automatically
forward messages they receive to the next listener before exposing them on their `C`
channel.

This approach removes the need to maintain an array of listeners belonging to a
channel and greatly simplifies the implementation.