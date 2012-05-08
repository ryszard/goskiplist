package skiplist

import "math/rand"

const p = 0.5
const maxLevel = 32 // maybe this should be configurable

type node struct {
	forward    []*node
	key, value interface{}
}

func (n *node) Next() *node {
	return n.forward[0]
}

func (n node) IsEnd() bool {
	return n.Next() == nil
//	return cap(n.forward) == 0
}

func (n *node) Key() interface{} {
	return n.key
}

func (n *node) Value() interface{} {
	return n.value
}

func (n node) Len() int {
	if n.IsEnd() {
		return 1
	}
	return 1 + n.Next().Len()
}

func (n node) level() int {
	return len(n.forward) - 1
}


type SkipList struct {
	lessThan    func(l, r interface{}) bool
	header *node
}

func (s SkipList) Len() int {
	// header shouldn't count as an element of the list.
	return s.header.Len() - 1
}

func (s SkipList) Iter() *node {
	return s.header.Next()
}

func (s SkipList) level() int {
	return s.header.level()
}

func (s SkipList) LessThan(l, r interface{}) bool {
	// nil is the maximum
	if l == nil {
		return false
	}
	return s.lessThan(l, r)
}

// Returns a new random level in 1..maxLevel.
func (s SkipList) randomLevel() (n int) {
	for n = 0; n < maxLevel && rand.Float64() < p; n++ {
	}
	return
}

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
		s.header.forward = s.header.forward[:s.level() - 1]
	}

	return candidate.Value(), true
}

func New(f func(l, r interface{}) bool) *SkipList {
	//end := &node{make([]*node, 0, 0), nil, nil}
	header := &node{[]*node{nil}, nil, nil}
	return &SkipList{lessThan: f, header: header}
}


type Comparable interface {
	LessThan(Comparable) bool
}

// NewComparable returns a SkipList that accepts skiplist.Comparable
// objects as keys.
func NewComparableMap() (s *SkipList) {
	comparator := func(left, right interface{}) bool {
		return left.(Comparable).LessThan(right.(Comparable))
	}
	return New(comparator)
	
}

// NewIntKey returns a SkipList that accepts int keys.
func NewIntMap() *SkipList {
	return New(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
}

func NewStringMap() *SkipList {
	return New(func(l, r interface{}) bool {
		return l.(string) < r.(string)
	})
}
