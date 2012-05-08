package skiplist

import "testing"
import "fmt"
import "math/rand"

func (s SkipList) printRepr() {

	for node := s.Front(); node != nil; node = node.Next() {
		fmt.Printf("%v: %v (level %d)\n", node.key, node.value, len(node.forward))
		for i, link := range node.forward {
			if link != nil {
				fmt.Printf("\t%d: -> %v\n", i, link.key)
			} else {
				fmt.Printf("\t%d: -> END\n", i)
			}
		}
	}
	fmt.Println()
}

func TestInitialization(t *testing.T) {
	s := NewMap(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
	if !s.lessThan(1, 2) {
		t.Errorf("Less than doesn't work correctly.")
	}
}

func TestHasNext(t *testing.T) {
	s := NewIntMap()
	s.Set(0, 0)
	node := s.header.Next()
	if node.Key() != 0 {
		t.Fatalf("We got the wrong node: %v.", node)
	}

	if node.HasNext() {
		t.Errorf("%v should be the last node.", node)
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
	s := NewIntMap()
	s.Set(0, 0)

	if value, present := s.Get(0); !(value == 0 && present) {
		t.Errorf("%v, %v instead of %v, %v", value, present, 0, true)
	}

	if value, present := s.Get(100); value != nil || present {
		t.Errorf("%v, %v instead of %v, %v", value, present, nil, false)
	}

}

func TestSet(t *testing.T) {
	s := NewIntMap()
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

func TestChange(t *testing.T) {
	s := NewIntMap()
	s.Set(0, 0)
	s.Set(1, 1)
	s.Set(2, 2)

	s.Set(0, 7)
	if value, _ := s.Get(0); value != 7 {
		t.Errorf("Value should be 7, not %d", value)
	}
	s.Set(1, 8)
	if value, _ := s.Get(1); value != 8 {
		t.Errorf("Value should be 8, not %d", value)
	}

}

func TestDelete(t *testing.T) {
	s := NewIntMap()
	for i := 0; i < 10; i++ {
		s.Set(i, i)
	}
	for i := 0; i < 10; i += 2 {
		s.Delete(i)
	}

	for i := 0; i < 10; i += 2 {
		if _, present := s.Get(i); present {
			t.Errorf("%d should not be present in s", i)
		}
	}
	if t.Failed() {
		s.printRepr()
	}

}

func TestLen(t *testing.T) {
	s := NewIntMap()
	for i := 0; i < 10; i++ {
		s.Set(i, i)
	}
	if length := s.Len(); length != 10 {
		t.Errorf("Length should be equal to 10, not %v.", length)
		s.printRepr()
	}
}

func TestIteration(t *testing.T) {
	s := NewIntMap()
	for i := 0; i < 20; i++ {
		s.Set(i, i)
	}

	seen := 0
	var lastKey int
	for i := s.Front(); i != nil; i = i.Next() {
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
	s := NewIntMap()
	insertions := [...]int{4, 1, 2, 9, 10, 7, 3}
	for _, i := range insertions {
		s.Set(i, i)
	}
	for _, i := range insertions {
		s.check(t, i, i)
	}

}

func makeRandomList(n int) (s *SkipList) {
	s = NewIntMap()
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
	s := NewIntMap()
	for i := 0; i < 10000; i++ {
		insert := rand.Int()
		s.Set(insert, insert)
	}
	var last int = 0
	for i := s.Front(); i != nil; i = i.Next() {
		if last != 0 && i.Key().(int) <= last {
			t.Errorf("Not in order!")
		}
		last = i.Key().(int)
	}
}

type MyComparable struct {
	value int
}

func (me MyComparable) LessThan(other Comparable) bool {
	return me.value < other.(MyComparable).value
}

func TestComparable(t *testing.T) {
	s := NewComparableMap()
	s.Set(MyComparable{0}, 0)
	s.Set(MyComparable{1}, 1)

	if val, _ := s.Get(MyComparable{0}); val != 0 {
		t.Errorf("Wrong value for MyComparable{0}. Should have been %d.", val)
	}
}

func TestNewStringMap(t *testing.T) {
	s := NewStringMap()
	s.Set("a", 1)
	s.Set("b", 2)
	if value, _ := s.Get("a"); value != 1 {
		t.Errorf("Expected 1, got %v.", value)
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
