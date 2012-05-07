package skiplist

import "math/rand"

const p = 0.5
const maxLevel = 32 // maybe this should be configurable

type SkipList struct {
	lessThan    func(l, r interface{}) bool
	header, end *node
}

func (s SkipList) Length() int {
	return s.header.Length()
}


func (n *node) next() *node {
	return n.forward[0]
}

func (n node) Length() int {
	if n.isEnd() {
		return -1
	}
	return 1 + n.forward[0].Length()
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
	return &node{make([]*node, level + 1, maxLevel), key, value}
}

// Returns a new random level in 1..maxLevel.
func (s SkipList) randomLevel() (n int) {
	for n = 0; n < maxLevel && rand.Float64() < p; n++ {
	}
	return
}

func (s *SkipList) Get(key interface{}) interface{} {
	current := s.header
	for i := s.level(); i >= 0; i-- {
		for s.LessThan(current.forward[i].key, key) {
			current = current.forward[i]
		}
	}
	if candidate := current.forward[0]; !candidate.isEnd() && candidate.key == key {
		return candidate.value
	}

	return nil
}

func (s *SkipList) Set(key, value interface{}) {
	update := make([]*node, maxLevel, maxLevel)
	current := s.header
	for i := s.level(); i >= 0; i-- {
		for s.LessThan(current.forward[i].key, key) {
			current = current.forward[i]
		}
		update[i] = current
	}

	// current is the last node whose key is less than key, update
	// is the node path to get from the header to current.

	if candidate := current.forward[0]; candidate.key == key {
		candidate.value = value
		return
	}

	newLevel := s.randomLevel()

	if currentLevel := s.level() ; newLevel > currentLevel {
		// there are no pointers for the higher levels in
		// update. Header should be there. Also add higher
		// level links to the header.
		for i := currentLevel + 1; i <= newLevel; i++ {
			update[i] = s.header
			s.header.forward = append(s.header.forward, s.end)
		}
	}

	newNode := makeNewNode(newLevel, key, value)
	
	for i := 0; i <= newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}
