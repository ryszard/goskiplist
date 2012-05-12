// Copyright 2012 Ric Szopa (Ryszard) <ryszard.szopa@gmail.com> All
// rights reserved.  Use of this source code is governed by a
// BSD-style license that can be found in the LICENSE file.

// Package skiplist implements skip list based maps and sets. For more
// information about skip lists take a look at
// http://en.wikipedia.org/wiki/Skip_list.

package skiplist

import "math/rand"

// TODO(ryszard):
//   - A separately seeded source of randomness
//   - Make maxLevel configurable
//   - Sets.

const p = 0.5
const maxLevel = 32 // maybe this should be configurable

type node struct {
	forward    []*node
	key, value interface{}
}

func (n *node) Next() *node {
	return n.forward[0]
}

func (n *node) HasNext() bool {
	return n.Next() != nil
}

func (n *node) Key() interface{} {
	return n.key
}

func (n *node) Value() interface{} {
	return n.value
}

type SkipList struct {
	lessThan func(l, r interface{}) bool
	header   *node
}

// Len returns the length of s.
func (s *SkipList) Len() (i int) {
	for node := s.Front(); node != nil; node = node.Next() {
		i++
	}
	return
}

// Front returns the first element of s in a way suitable for
// iteration.
func (s *SkipList) Front() *node {
	return s.header.Next()
}

func (s *SkipList) level() int {
	return len(s.header.forward) - 1
}

func (s *SkipList) LessThan(l, r interface{}) bool {
	// nil is the maximum
	if l == nil {
		return false
	}
	return s.lessThan(l, r)
}

// Returns a new random level.
func (s SkipList) randomLevel() (n int) {
	for n = 0; n < maxLevel && rand.Float64() < p; n++ {
	}
	return
}

// Get returns the value associated with key from s (nil if the key is
// not present in s). The second return value is true when the key is
// present.
func (s *SkipList) Get(key interface{}) (value interface{}, present bool) {
	candidate := s.getPath(nil, key)

	if candidate != nil && candidate.key == key {
		return candidate.value, true
	}

	return nil, false
}

// getPath populates update with nodes that constitute the path to the
// node that may contain key. The candidate node will be returned. If
// update is nil, it will be left alone (the candidate node will still
// be returned). If update is not nil, but it doesn't have enough
// slots for all the nodes in the path, getPath will panic.
func (s *SkipList) getPath(update []*node, key interface{}) *node {
	current := s.header
	for i := s.level(); i >= 0; i-- {
		for current.forward[i] != nil && s.LessThan(current.forward[i].key, key) {
			current = current.forward[i]
		}
		if update != nil {
			update[i] = current
		}
	}
	return current.Next()
}

// Sets set the value associated with key in s.
func (s *SkipList) Set(key, value interface{}) {

	// s.level starts from 0, so we need to allocate one 
	update := make([]*node, s.level()+1, maxLevel)
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

	newNode := &node{make([]*node, newLevel+1, maxLevel), key, value}

	for i := 0; i <= newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}

// Delete removes node with key key from s and returns its value. A
// second value, a boolean, is return to indicate if key was present
// in s.
func (s *SkipList) Delete(key interface{}) (value interface{}, present bool) {
	update := make([]*node, s.level()+1, maxLevel)
	candidate := s.getPath(update, key)

	if candidate.key != key {
		return nil, false
	}

	for i := 0; i <= s.level() && update[i].forward[i] == candidate; i++ {
		update[i].forward[i] = candidate.forward[i]
	}

	for s.level() > 0 && s.header.forward[s.level()] == nil {
		s.header.forward = s.header.forward[:s.level()-1]
	}

	return candidate.Value(), true
}

// New returns a new SkipList that will use lessThan as the comparison
// function. lessThan should be linear order on keys you intend to use
// with the SkipList.
func NewMap(lessThan func(l, r interface{}) bool) *SkipList {
	return &SkipList{lessThan, &node{[]*node{nil}, nil, nil}}
}

type Ordered interface {
	LessThan(Ordered) bool
}

// NewOrderedMap returns a SkipList that accepts skiplist.Ordered
// objects as keys.
func NewOrderedMap() (s *SkipList) {
	comparator := func(left, right interface{}) bool {
		return left.(Ordered).LessThan(right.(Ordered))
	}
	return NewMap(comparator)

}

// NewIntKey returns a SkipList that accepts int keys.
func NewIntMap() *SkipList {
	return NewMap(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
}

// NewStringMap returns a SkipList accepting strings as keys.
func NewStringMap() *SkipList {
	return NewMap(func(l, r interface{}) bool {
		return l.(string) < r.(string)
	})
}
