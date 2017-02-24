// redis like sorted set
package skiplist

import "math"

type ZSet struct {
	key2Score map[interface{}]*zsetScore
	sl        *SkipList
	pool      *zsetScorePool
}

type zsetScore struct {
	score   interface{}
	counter int64
}

type zsetScorePool struct {
	pool    chan *zsetScore
	counter int64
}

func newzsetScorePool(cap int) *zsetScorePool {
	return &zsetScorePool{
		pool: make(chan *zsetScore, cap),
	}
}

func (p *zsetScorePool) Get(score interface{}) *zsetScore {
	select {
	case s := <-p.pool:
		s.score = score
		p.counter++
		s.counter = p.counter
		return s
	default:
		p.counter++
		return &zsetScore{
			score:   score,
			counter: p.counter,
		}
	}
}

func (p *zsetScorePool) Put(s *zsetScore) {
	select {
	case p.pool <- s:
	default:
	}
}

func NewCustomZSet(scoreLessThan func(l, r interface{}) bool) *ZSet {
	return &ZSet{
		key2Score: make(map[interface{}]*zsetScore),
		sl: NewCustomMap(func(l, r interface{}) bool {
			lzs := l.(*zsetScore)
			rzs := r.(*zsetScore)
			if scoreLessThan(lzs.score, rzs.score) {
				return true
			} else if lzs.score == rzs.score && lzs.counter < rzs.counter {
				return true
			} else {
				return false
			}
		}),
		pool: newzsetScorePool(128),
	}
}

func NewZSet() *ZSet {
	return NewCustomZSet(func(l, r interface{}) bool {
		return l.(Ordered).LessThan(r.(Ordered))
	})
}

func (z *ZSet) Add(key interface{}, score interface{}) bool {
	curZScore, ok := z.key2Score[key]
	if ok {
		if score != curZScore.score { // update
			z.sl.Delete(curZScore)
			z.pool.Put(curZScore)
			zScore := z.pool.Get(score)
			z.sl.Set(zScore, key)
			z.key2Score[key] = zScore
		}
	} else {
		zScore := z.pool.Get(score)
		z.key2Score[key] = zScore
		z.sl.Set(zScore, key)
	}
	return true
}

func (z *ZSet) Update(key interface{}, score interface{}) bool {
	curZScore, ok := z.key2Score[key]
	if !ok {
		return false
	}
	if score != curZScore.score { // update
		z.sl.Delete(curZScore)
		z.pool.Put(curZScore)
		zScore := z.pool.Get(score)
		z.sl.Set(zScore, key)
		z.key2Score[key] = zScore
	}
	return true
}

func (z *ZSet) Remove(key interface{}) bool {
	curZScore, ok := z.key2Score[key]
	if !ok {
		return false
	}
	z.sl.Delete(curZScore)
	z.pool.Put(curZScore)
	delete(z.key2Score, key)
	return true
}

func (z *ZSet) Rank(key interface{}) uint32 {
	curZScore, ok := z.key2Score[key]
	if !ok {
		return 0
	}
	return z.sl.Rank(curZScore)
}

func (z *ZSet) Score(key interface{}) interface{} {
	curZScore, _ := z.key2Score[key]
	return curZScore.score
}

func (z *ZSet) RangeByRank(rankFrom uint32, rankTo uint32) [][2]interface{} { // [rankFrom, rankTo]
	if rankTo > uint32(z.sl.Len()) {
		rankTo = uint32(z.sl.Len())
	}

	if rankTo < rankFrom {
		return nil
	}

	iter := z.sl.GetElemByRank(rankFrom)
	if iter == nil {
		return nil
	}
	keys := make([][2]interface{}, 0, int(rankTo-rankFrom+1))
	for i := rankFrom; i <= rankTo; i++ {
		keys = append(keys, [2]interface{}{iter.Value(), iter.Key().(*zsetScore).score})
		if !iter.Next() {
			break
		}
	}
	return keys
}

func (z *ZSet) RangeByScore(scoreFrom interface{}, scoreTo interface{}) []interface{} { // [scoreFrom, scoreTo]
	iter := z.sl.Range(&zsetScore{score: scoreFrom}, &zsetScore{score: scoreTo, counter: math.MaxInt64})
	keys := make([]interface{}, 0, 8)
	rangeIter := iter.(*rangeIterator)
	for rangeIter.Next() {
		keys = append(keys, rangeIter.Value())
	}
	return keys
}

func (z *ZSet) Card() int { // 集合元素个数
	return len(z.key2Score)
}

func (z *ZSet) Foreach(fn func(key interface{}, score interface{})) {
	iter := z.sl.Iterator()
	for iter.Next() {
		fn(iter.Value(), iter.Key().(*zsetScore).score)
	}
}

func (z *ZSet) Clear() {
	z.key2Score = make(map[interface{}]*zsetScore)
	z.sl.Clear()
}

func (z *ZSet) Marshal() [][2]interface{} {
	elements := make([][2]interface{}, 0, len(z.key2Score))
	iter := z.sl.Iterator()
	for iter.Next() {
		elements = append(elements, [2]interface{}{iter.Value(), iter.Key().(*zsetScore).score})
	}
	return elements
}

func (z *ZSet) Unmarshal(elements [][2]interface{}) bool {
	for i, elem := range elements {
		zScore := z.pool.Get(elem[1])
		z.key2Score[elem[0]] = zScore
		elements[i][0] = zScore
		elements[i][1] = elem[0]
	}
	return z.sl.FillBySortedSlice(elements)
}
