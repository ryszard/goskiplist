package skiplist

import "testing"
import "fmt"
import "math/rand"

func (s SkipList) printRepr() {

	for node := s.header; !node.IsEnd(); node = node.forward[0] {
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

func TestIsEnd(t *testing.T) {
	s := NewIntKey()
	if !s.end.IsEnd() {
	 	t.Errorf("IsEnd() is flase for s.end.")
	}

	if s.header.IsEnd() {
		t.Errorf("IsEnd() is true for s.header.")
	}

	s.Set(0, 0)
	node := s.header.Next()
	if node.Key() != 0 {
		t.Fatalf("We got the wrong node: %v.", node)
	}

	if node.IsEnd() {
		t.Errorf("IsEnd() should be false for %v.", node)
	}

	if node == s.end {
		t.Errorf("%v should not be equal to s.end.", node)
	}

	if node.Next() != s.end {
		t.Errorf("node.next should not be equal to s.end (was %v).", node, node.Next())
	}

}

func (s SkipList) check(t *testing.T, key, wanted int) bool {
	if got, _ := s.Get(key); got != wanted {
		t.Errorf("Wanted %v, got %v.", wanted, got)
		return true
	}
	return false
}

func TestGet(t *testing.T) {
	s := NewIntKey()
	s.Set(0, 0)
	
	if value, present := s.Get(0); !(value == 0 && present) {
		t.Errorf("%v, %v instead of %v, %v", value, present, 0, true)
	}
	
	if value, present := s.Get(100); value != nil || present {
		t.Errorf("%v, %v instead of %v, %v", value, present, nil, false)
	}


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
		t.Errorf("%v", s.header.Next())
	}
	s.check(t, 1, 1)

}

func TestDelete(t *testing.T) {
	s := NewIntKey()
	for i := 0; i < 10; i++ {
		s.Set(i, i)
	}
	for i := 0; i < 10; i+=2 {
		s.Delete(i)
	}

	for i := 0; i < 10; i+=2 {
		if _, present := s.Get(i); present {
			t.Errorf("%d should not be present in s", i)
		}
	}
	if t.Failed() {
		s.printRepr()
	}
	
}

func TestLen(t *testing.T) {
	s := NewIntKey()
	for i := 0; i < 10; i++ {
		s.Set(i, i)
	}
	if length := s.Len(); length != 10 {
		t.Errorf("Length should be equal to 10, not %v.", length)
		s.printRepr()
	}
}

func TestIterator(t *testing.T) {
	s := NewIntKey()
	for i := 0; i < 20; i++ {
		s.Set(i, i)
	}

	seen := 0
	var lastKey int
	for i := s.Iter(); !i.IsEnd(); i = i.Next() {
		seen++
		lastKey = i.Key().(int)
		if i.Key() != i.Value() {
			t.Errorf("Wrong value for key %v: %v.", i.Key(), i.Value())
		}
	}

	if seen != s.Len() {
		t.Errorf("Not all the items in s where iterated through (seen %d, should have seen %d). Last one seen was %d.", seen, s.Len(), lastKey)
	}
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

}

func makeRandomList(n int) (s *SkipList) {
	s = NewIntKey()
	for i := 0; i < n; i++ {
		insert := rand.Int()
		s.Set(insert, insert)
	}
	return
}


func LookupBenchmark(b *testing.B, n int) {
	b.StopTimer()
	s := makeRandomList(n)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Get(rand.Int())
	}
}


// Make sure that all the keys are unique and are returned in order.
func TestSanity(t *testing.T) {
	s := NewIntKey()
	for i := 0; i < 10000; i++ {
		insert := rand.Int()
		s.Set(insert, insert)
	}
	var last int = 0
	for i := s.Iter(); !i.IsEnd(); i = i.Next() {
		if last != 0 && i.Key().(int) <= last {
			t.Errorf("Not in order!")
		}
		last = i.Key().(int)
	}
}


func BenchmarkLookup16(b *testing.B) {
	LookupBenchmark(b, 16)
}


func BenchmarkLookup256(b *testing.B) {
	LookupBenchmark(b, 256)
}


func BenchmarkLookup65536(b *testing.B) {
	LookupBenchmark(b, 65536)
}