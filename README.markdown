About
=====

This is a library implementing skip lists for the Go programming
language (http://golang.org/).

Skip lists are a data structure that can be used in place of
balanced trees. Skip lists use probabilistic balancing rather than
strictly enforced balancing and as a result the algorithms for
insertion and deletion in skip lists are much simpler and
significantly faster than equivalent algorithms for balanced trees.

Skip lists were first described in
[Pugh, William (June 1990)](ftp://ftp.cs.umd.edu/pub/skipLists/skiplists.pdf). "Skip
lists: a probabilistic alternative to balanced trees". Communications
of the ACM 33 (6): 668–676

Installing
==========

    $ go get github.com/ryszard/goskiplist/skiplist
	
Example
=======

```go
package main

import (
	"github.com/ryszard/goskiplist/skiplist"
	"fmt"
)

func main() {
	s := skiplist.NewIntMap()
	s.Set(7, "seven")
	s.Set(1, "one")
	s.Set(0, "zero")
	s.Set(5, "five")
	s.Set(9, "nine")
	s.Set(10, "ten")
	s.Set(3, "three")

	value, ok := s.Get(0)
	if ok {
		fmt.Println(value)
	}
	// prints: 
	//	zero


	s.Delete(7)

	value, ok := s.Get(7)
	if ok {
		fmt.Println(value)
	}
	// prints: nothing.

	s.Set(9, "niner")

	// Iterate through all the elements, in order.
	for i := s.Iterator(); i.Next(); {
		fmt.Printf("%d: %s\n", i.Key(), i.Value())
	}
	// prints: 
	//	0: zero
	// 	1: one
	// 	3: three
	// 	5: five
	// 	9: niner
	// 	10: ten

	// Iterate only through elements in some range.
	for i := s.Range(3, 10); i.Next(); {
		fmt.Printf("%d: %s\n", i.Key(), i.Value())
	}
	// prints: 
	// 	3: three
	// 	5: five
	// 	9: niner

}
```

Full documentation
==================

Read it [online](http://go.pkgdoc.org/github.com/ryszard/goskiplist/skiplist) or run 

    $ go doc github.com/ryszard/goskiplist/skiplist
