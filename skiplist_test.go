package skiplist

import "testing"
import "fmt"
import "math/rand"

func (s SkipList) printRepr() {

	for node := s.header; !node.isEnd(); node = node.forward[0] {
		fmt.Printf("%v: %v (level %d)\n", node.key, node.value, node.level())
		for i, link := range node.forward {
			fmt.Printf("\t%d: -> %v\n", i, link.key)
		}
	}
	fmt.Println()
}

func TestInitialization(t *testing.T) {
	s := New(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
	if !s.lessThan(1, 2) {
		t.Errorf("Less than doesn't work correctly.")
	}
}

func (s SkipList) check(t *testing.T, key, wanted int) bool {
	if got := s.Get(key); got != wanted {
		t.Errorf("Wanted %v, got %v.", wanted, got)
		return true
	}
	return false
}

func TestSet(t *testing.T) {
	s := NewIntKey()
	if l := s.Len(); l != 0 {
		t.Errorf("Len is not 0, it is %v", l)
	}

	s.Set(0, 0)
	s.Set(1, 1)
	if l := s.Len(); l != 2 {
		t.Errorf("Len is not 2, it is %v", l)
	}
	if s.check(t, 0, 0) {
		t.Errorf("%v", s.header.next())
	}
	s.check(t, 1, 1)

}

func TestSomeMore(t *testing.T) {
	s := NewIntKey()
	insertions := [...]int{4, 1, 2, 9, 10, 7, 3}
	for _, i := range insertions {
		s.Set(i, i)
	}
	for _, i := range insertions {
		s.check(t, i, i)
	}
	//s.printRepr()

}

func BenchmarkInsertion(b *testing.B) {
	s := NewIntKey()
	for i := 0; i < b.N; i++ {
		insert := rand.Int()
		s.Set(insert, insert)
	}
}
