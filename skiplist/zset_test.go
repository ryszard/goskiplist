package skiplist

import (
	"math/rand"
	"testing"
)

func TestZSet(t *testing.T) {
	zs := NewCustomZSet(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
	for i := 0; i < 100; i++ {
		zs.Add(i, i*10)
	}
	if zs.Card() != 100 {
		t.Errorf("after add 100, zset length should be 100")
	}

	for i := 0; i < 100; i++ {
		if zs.Rank(i) != uint32(i+1) {
			t.Errorf("rank error")
		}
	}

	for i, ks := range zs.RangeByRank(1, 10000) {
		if ks[1].(int) != i*10 || ks[0].(int) != i {
			t.Errorf("rangebyrank error")
		}
	}

	for i, k := range zs.RangeByScore(0, 1000) {
		if k.(int) != i {
			t.Errorf("rangbyscore error")
		}
	}

	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			ok := zs.Remove(i)
			if !ok {
				t.Errorf("remove failed")
			}
		}
	}
	if zs.Card() != 50 {
		t.Errorf("after remove 50, zset length should be 50")
	}

	zsSlice := zs.Marshal()
	for i, elem := range zsSlice {
		if elem[0].(int) != (i*2+1) || elem[1].(int) != 10*(i*2+1) {
			t.Errorf("marshal error")
		}
	}

	zs.Clear()

	zs.Unmarshal(zsSlice)
	for i := 0; i < 100; i++ {
		if i%2 != 0 {
			if zs.Rank(i) != uint32(i/2)+1 || zs.Score(i).(int) != i*10 {
				t.Errorf("unmarshal error")
			}
		}
	}
}

func TestZSetRank(t *testing.T) {
	zs := NewCustomZSet(func(l, r interface{}) bool {
		return l.(int) > r.(int)
	})
	zs.Add("foo", 12)
	zs.Add("bar", 12)
	if zs.Rank("foo") != 1 || zs.Rank("bar") != 2 {
		t.Errorf("rank perform wrong")
	}
	zs.Add("bar", 13)
	if zs.Rank("bar") != 1 || zs.Rank("foo") != 2 {
		t.Errorf("rank perform wrong")
	}
}

func shuffleArray(array []int) {
	for len(array) != 0 {
		pos := rand.Intn(len(array))
		array[0], array[pos] = array[pos], array[0]
		array = array[1:]
	}
}

type customType int

func (c customType) LessThan(o Ordered) bool {
	return c < o.(customType)
}

func TestZSet2(t *testing.T) {
	zs := NewCustomZSet(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
	length := 1000000
	array := make([]int, length)
	for i := 0; i < length; i++ {
		array[i] = i
	}
	shuffleArray(array)
	for _, v := range array {
		zs.Add(v, v)
	}
	for _, v := range array {
		if zs.Rank(v) != uint32(v+1) {
			t.Fatalf("rank perform wrong")
		}
	}

	rankFrom := uint32(rand.Intn(len(array))) + 1
	for i, ks := range zs.RangeByRank(rankFrom, uint32(len(array))) {
		if uint32(ks[0].(int)+1) != uint32(i)+rankFrom {
			t.Fatalf("range by rank perform wrong")
		}
	}

	zsSlice := zs.Marshal()
	zs.Clear()
	zs.Unmarshal(zsSlice)

	for _, v := range array {
		if zs.Rank(v) != uint32(v+1) {
			t.Fatalf("rank perform wrong")
		}
	}

	for _, v := range array {
		zs.Update(v, -v)
	}

	zs.Foreach(func(key interface{}, score interface{}) {
		if key.(int) != -(score.(int)) {
			t.Fatalf("foreach perform wrong")
		}
	})

	for _, v := range array {
		zs.Remove(v)
	}

	if zs.RangeByRank(100, 300) != nil {
		t.Fatalf("range by rank perform wrong")
	}

	if zs.Rank(1) != 0 {
		t.Fatalf("rank perform wrong")
	}

	if zs.Update(1, 99) {
		t.Fatalf("update perform wrong")
	}

	if zs.Remove(1) {
		t.Fatalf("remove perform wrong")
	}

	zs2 := NewZSet()
	zs2.Add("foo", customType(1))
	zs2.Add("bar", customType(2))
	if zs2.Rank("foo") != 1 || zs2.Rank("bar") != 2 {
		t.Fatalf("NewZSet perform wrong")
	}
}

func BenchmarkZSetAdd(b *testing.B) {
	zs := NewCustomZSet(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})

	length := 10000000
	array := make([]int, length)
	for i := 0; i < length; i++ {
		array[i] = i
	}
	shuffleArray(array)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zs.Add(i, array[i])
	}
}

func BenchmarkZSetRank(b *testing.B) {
	zs := NewCustomZSet(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})

	array := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		array[i] = i
	}
	shuffleArray(array)
	for i := 0; i < b.N; i++ {
		zs.Add(array[i], array[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if zs.Rank(i) != uint32(i+1) {
			b.Fatalf("rank perform wrong")
		}
	}
}
