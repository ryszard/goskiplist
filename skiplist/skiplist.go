// Copyright 2012 Google Inc. All rights reserved.
// Author: Ric Szopa (Ryszard) <ryszard.szopa@gmail.com>

// Package skiplist implements skip list based maps and sets.
//
// Skip lists are a data structure that can be used in place of
// balanced trees. Skip lists use probabilistic balancing rather than
// strictly enforced balancing and as a result the algorithms for
// insertion and deletion in skip lists are much simpler and
// significantly faster than equivalent algorithms for balanced trees.
//
// Skip lists were first described in Pugh, William (June 1990). "Skip
// lists: a probabilistic alternative to balanced
// trees". Communications of the ACM 33 (6): 668–676
package skiplist

import "math/rand"

// TODO(ryszard):
//   - A separately seeded source of randomness

// p is the fraction of nodes with level i pointers that also have
// level i+1 pointers. p equal to 1/4 is a good value from the point
// of view of speed and space requirements. If variability of running
// times is a concern, 1/2 is a better value for p.
const p = 0.25

const DefaultMaxLevel = 32

// A node is a container for key-value pairs that are stored in a skip
// list.
type node struct {
	forward    []*node
	key, value interface{}
}

// next returns the next node in the skip list containing n.
func (n *node) next() *node {
	if len(n.forward) == 0 {
		return nil
	}
	return n.forward[0]
}

// hasNext returns true if n has a next node.
func (n *node) hasNext() bool {
	return n.next() != nil
}

// A SkipList is a map-like data structure that maintains an ordered
// collection of key-value pairs. Insertion, lookup, and deletion are
// all O(log n) operations. A SkipList can efficiently store up to
// 2^MaxLevel items.
//
// To iterate over a skip list (where s is a
// *SkipList):
//
//	for i := s.Iterator(); i.Next(); {
//		// do something with i.Key() and i.Value()
//	}
type SkipList struct {
	lessThan func(l, r interface{}) bool
	header   *node
	length   int
	// MaxLevel determines how many items the SkipList can store
	// efficiently (2^MaxLevel).
	//
	// It is safe to increase MaxLevel to accomodate more
	// elements. If you decrease MaxLevel and the skip list
	// already contains nodes on higer levels, the effective
	// MaxLevel will be the greater of the new MaxLevel and the
	// level of the highest node.
	//
	// A SkipList with MaxLevel equal to 0 is equivalent to a
	// standard linked list and will not have any of the nice
	// properties of skip lists (probably not what you want).
	MaxLevel int
}

// Len returns the length of s.
func (s *SkipList) Len() int {
	return s.length
}

// Iterator is an interface that you can use to iterate through the
// skip list (in its entirety or fragments). For an use example, see
// the documentation of SkipList.
//
// Key and Value return the key and the value of the current node.
type Iterator interface {
	// Next returns true if the iterator contains more elements
	// and andvances its state to the next element if that is
	// possible.
	Next() bool
	// Key returns the current key.
	Key() interface{}
	// Value returns the current value.
	Value() interface{}
}

type iter struct {
	key, value interface{}
	current    *node
}

func (i iter) Key() interface{} {
	return i.key
}

func (i iter) Value() interface{} {
	return i.value
}

func (i *iter) Next() bool {
	if i.current.hasNext() {
		i.current = i.current.next()
		i.key = i.current.key
		i.value = i.current.value
		return true
	}
	return false
}

type rangeIterator struct {
	iter
	limit    interface{}
	skipList *SkipList
}

func (i *rangeIterator) Next() bool {
	if i.current.hasNext() {
		next := i.current.next()
		if !i.skipList.lessThan(next.key, i.limit) {
			return false
		}
		i.current = i.current.next()
		i.key = i.current.key
		i.value = i.current.value
		return true
	}
	return false
}

// Iterator returns an Iterator that will go through all elements s.
func (s *SkipList) Iterator() Iterator {
	return &iter{current: s.header}
}

// Range returns an iterator that will go through all the
// elements of the skip list that are greater or equal than from, but
// less than to.
func (s *SkipList) Range(from, to interface{}) Iterator {
	start := s.getPath(nil, from)
	return &rangeIterator{iter: iter{current: &node{forward: []*node{start}}}, limit: to, skipList: s}
}

func (s *SkipList) level() int {
	return len(s.header.forward) - 1
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (s *SkipList) effectiveMaxLevel() int {
	return maxInt(s.level(), s.MaxLevel)
}

// Returns a new random level.
func (s SkipList) randomLevel() (n int) {
	for n = 0; n < s.effectiveMaxLevel() && rand.Float64() < p; n++ {
	}
	return
}

// Get returns the value associated with key from s (nil if the key is
// not present in s). The second return value is true when the key is
// present.
func (s *SkipList) Get(key interface{}) (value interface{}, ok bool) {
	candidate := s.getPath(nil, key)

	if candidate == nil || candidate.key != key {
		return nil, false
	}

	return candidate.value, true
}

// GetGreaterOrEqual finds the node whose key is greater than or equal
// to min. It returns its value, its actual key, and whether such a
// node is present in the skip list.
func (s *SkipList) GetGreaterOrEqual(min interface{}) (actualKey, value interface{}, ok bool) {
	candidate := s.getPath(nil, min)

	if candidate != nil {
		return candidate.value, candidate.key, true
	}
	return nil, nil, false
}

// getPath populates update with nodes that constitute the path to the
// node that may contain key. The candidate node will be returned. If
// update is nil, it will be left alone (the candidate node will still
// be returned). If update is not nil, but it doesn't have enough
// slots for all the nodes in the path, getPath will panic.
func (s *SkipList) getPath(update []*node, key interface{}) *node {
	current := s.header
	for i := s.level(); i >= 0; i-- {
		for current.forward[i] != nil && s.lessThan(current.forward[i].key, key) {
			current = current.forward[i]
		}
		if update != nil {
			update[i] = current
		}
	}
	return current.next()
}

// Sets set the value associated with key in s.
func (s *SkipList) Set(key, value interface{}) {
	if key == nil {
		panic("nil keys are not supported")
	}
	// s.level starts from 0, so we need to allocate one.
	update := make([]*node, s.level()+1, s.effectiveMaxLevel()+1)
	candidate := s.getPath(update, key)

	if candidate != nil && candidate.key == key {
		candidate.value = value
		return
	}

	newLevel := s.randomLevel()

	if currentLevel := s.level(); newLevel > currentLevel {
		// there are no pointers for the higher levels in
		// update. Header should be there. Also add higher
		// level links to the header.
		for i := currentLevel + 1; i <= newLevel; i++ {
			update = append(update, s.header)
			s.header.forward = append(s.header.forward, nil)
		}
	}

	newNode := &node{make([]*node, newLevel+1, s.effectiveMaxLevel()+1), key, value}

	for i := 0; i <= newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
	s.length++
}

// Delete removes the node with the given key.
//
// It returns the old value and whether the node was present.
func (s *SkipList) Delete(key interface{}) (value interface{}, ok bool) {
	update := make([]*node, s.level()+1, s.effectiveMaxLevel())
	candidate := s.getPath(update, key)

	if candidate == nil || candidate.key != key {
		return nil, false
	}

	for i := 0; i <= s.level() && update[i].forward[i] == candidate; i++ {
		update[i].forward[i] = candidate.forward[i]
	}

	for s.level() > 0 && s.header.forward[s.level()] == nil {
		s.header.forward = s.header.forward[:s.level()]
	}
	s.length--
	return candidate.value, true
}

// NewCustomMap returns a new SkipList that will use lessThan as the
// comparison function. lessThan should define a linear order on keys
// you intend to use with the SkipList.
func NewCustomMap(lessThan func(l, r interface{}) bool) *SkipList {
	return &SkipList{
		lessThan: lessThan,
		header:   &node{forward: []*node{nil}},
		MaxLevel: DefaultMaxLevel,
	}
}

// Ordered is an interface which can be linearly ordered by the
// LessThan method.
type Ordered interface {
	LessThan(Ordered) bool
}

// New returns a new SkipList.
//
// Its keys must implement the Ordered interface.
func New() *SkipList {
	comparator := func(left, right interface{}) bool {
		return left.(Ordered).LessThan(right.(Ordered))
	}
	return NewCustomMap(comparator)

}

// NewIntKey returns a SkipList that accepts int keys.
func NewIntMap() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
}

// NewStringMap returns a SkipList that accepts string keys.
func NewStringMap() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(string) < r.(string)
	})
}

// Set is an ordered set data structure.
//
// Its elements must implement the Ordered interface. It uses a
// SkipList for storage, and it gives you similar performance
// guarantees.
//
// To iterate over a set (where s is a
// *Set):
//
//	for i := s.Iterator(); i.Next(); {
//		// do something with i.Key().
//		// i.Value() will be nil.
//	}
type Set struct {
	skiplist SkipList
}

// NewSet returns a new Set.
func NewSet() *Set {
	comparator := func(left, right interface{}) bool {
		return left.(Ordered).LessThan(right.(Ordered))
	}
	return NewCustomSet(comparator)
}

// NewCustomSet returns a new Set that will use lessThan as the
// comparison function. lessThan should define a linear order on
// elements you intend to use with the Set.
func NewCustomSet(lessThan func(l, r interface{}) bool) *Set {
	return &Set{skiplist: SkipList{
		lessThan: lessThan,
		header:   &node{forward: []*node{nil}},
		MaxLevel: DefaultMaxLevel,
	}}
}

// NewIntSet returns a new Set that accepts int elements.
func NewIntSet() *Set {
	return NewCustomSet(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
}

// NewStringSet returns a new Set that acceots string elements.
func NewStringSet() *Set {
	return NewCustomSet(func(l, r interface{}) bool {
		return l.(string) < r.(string)
	})
}

// Add adds key to s.
func (s *Set) Add(key interface{}) {
	s.skiplist.Set(key, nil)
}

// Remove tries to remove key from the set. It returns true if key was
// present.
func (s *Set) Remove(key interface{}) (ok bool) {
	_, ok = s.skiplist.Delete(key)
	return ok
}

// Len returns the length of the set.
func (s *Set) Len() int {
	return s.skiplist.Len()
}
// Contains returns true if key is present in s.
func (s *Set) Contains(key interface{}) bool {
	_, ok := s.skiplist.Get(key)
	return ok
}

func (s *Set) Iterator() Iterator {
	return s.skiplist.Iterator()
}

// Range returns an iterator that will go through all the elements of
// the set that are greater or equal than from, but less than to.
func (s *Set) Range(from, to interface{}) Iterator {
	return s.skiplist.Range(from, to)
}

// SetMaxLevel sets MaxLevel in the underlying skip list.
func (s *Set) SetMaxLevel(newMaxLevel int) {
	s.skiplist.MaxLevel = newMaxLevel
}

// GetMaxLevel returns MaxLevel fo the underlying skip list.
func (s *Set) GetMaxLevel() int {
	return s.skiplist.MaxLevel
}
