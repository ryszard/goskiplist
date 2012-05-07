package skiplist

import "math/rand"

const p = 0.5
const maxLevel = 32 // maybe this should be configurable

type SkipList struct {
	lessThan    func(l, r interface{}) bool
	header, end *node
}

func (s SkipList) Len() int {
	// header and and end do not count as elements of the list.
	return s.header.Len() - 1
}

func (n *node) Next() *node {
	return n.forward[0]
}

func (s SkipList) Iter() *node {
	return s.header.Next()
}

func (n *node) Key() interface{} {
	return n.key
}

func (n *node) Value() interface{} {
	return n.value
}

func (n node) Len() int {
	if n.isEnd() {
		return 0
	}
	return 1 + n.Next().Len()
}

func (n node) level() int {
	return len(n.forward) - 1
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

type node struct {
	forward    []*node
	key, value interface{}
}

func (n node) isEnd() bool {
	return cap(n.forward) == 0
}

func (n node) HasNext() bool {
	return !n.isEnd()
}

func New(f func(l, r interface{}) bool) *SkipList {
	end := &node{make([]*node, 0, 0), nil, nil}
	header := &node{[]*node{end}, nil, nil}
	return &SkipList{lessThan: f, header: header, end: end}
}

func NewIntKey() *SkipList {
	return New(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
}

func makeNewNode(level int, key, value interface{}) *node {
	return &node{make([]*node, level+1, maxLevel), key, value}
}

// Returns a new random level in 1..maxLevel.
func (s SkipList) randomLevel() (n int) {
	for n = 0; n < maxLevel && rand.Float64() < p; n++ {
	}
	return
}

func (s *SkipList) Get(key interface{}) interface{} {
	candidate := s.getPath(nil, key)

	if !candidate.isEnd() && candidate.key == key {
		return candidate.value
	}

	return nil
}

// getPath populates update with nodes that constitute the path to the
// node that may contain key. The candidate node will be returned. If
// update is nil, it will be left alone (the candidate node will still
// be returned). If update is not nil, but it doesn't have enough
// slots for all the nodes in the path, getPath will panic.
func (s *SkipList) getPath(update []*node, key interface{}) *node {
	current := s.header
	for i := s.level(); i >= 0; i-- {
		for s.LessThan(current.forward[i].key, key) {
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

	if candidate.key == key {
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
			s.header.forward = append(s.header.forward, s.end)
		}
	}

	newNode := makeNewNode(newLevel, key, value)

	for i := 0; i <= newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}

func (s *SkipList) Delete(key interface{}) (interface{}, bool) {
	update := make([]*node, s.level()+1, maxLevel)
	candidate := s.getPath(update, key)

	if candidate.key != key {
		return nil, false
	}

	for i := 0; i <= s.level() && update[i].forward[i] == candidate; i++ {
		update[i].forward[i] = candidate.forward[i]
	}

	for s.level() > 0 && s.header.forward[s.level()] == s.end {
		s.header.forward = s.header.forward[:s.level() - 1]
	}

	return candidate.Value(), true
}

// TODO(szopa): deletion, test that there are no duplicates, test that
// the values are in order, get a unique source of randomness (that
// can be seeded separately).
