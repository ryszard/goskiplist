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

Full documentation
==================

Read it [online](http://go.pkgdoc.org/github.com/ryszard/goskiplist/skiplist) or run 

    $ go doc github.com/ryszard/goskiplist/skiplist
