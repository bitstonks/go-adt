package set

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type E = int

func contents(s Set[E]) []E {
	res := make([]E, 0, len(s))
	for k := range s {
		res = append(res, k)
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return res
}

func TestNew(t *testing.T) {
	t.Parallel()

	empty := []E{}
	one := []E{1}
	five := []E{0, 1, 2, 3, 4}

	check := func(t *testing.T, expected []E, set Set[E]) {
		assert.Equal(t, expected, contents(set))
	}

	t.Run("none", func(t *testing.T) { check(t, empty, New[E]()) })
	t.Run("{}", func(t *testing.T) { check(t, empty, New[E](empty...)) })
	t.Run("{1}", func(t *testing.T) { check(t, one, New[E](one...)) })
	t.Run("{0,1,2,3,4}", func(t *testing.T) { check(t, five, New[E](five...)) })
}

var snil Set[E]
var null = New[E]()
var s1 = New[E](0, 1, 2, 3, 4)
var s2 = New[E](3, 4, 5, 6, 7)
var s3 = New[E](6, 7, 8, 9, 10)
var subset3 = New[E](6, 7)

func TestCopy(t *testing.T) {
	t.Parallel()

	check := func(t *testing.T, expected, s Set[E]) {
		assert.Equal(t, expected, s.Copy())
	}
	t.Run("nil", func(t *testing.T) { check(t, null, nil) })
	t.Run("null", func(t *testing.T) { check(t, null, null) })
	t.Run("s1", func(t *testing.T) { check(t, s1, s1) })
	t.Run("s2", func(t *testing.T) { check(t, s2, s2) })
	t.Run("s3", func(t *testing.T) { check(t, s3, s3) })
}

func TestContains(t *testing.T) {
	t.Parallel()

	check := func(t *testing.T, expected bool, s Set[E], key E, keys ...E) {
		assert.Equal(t, expected, s.Contains(key, keys...))
	}

	t.Run("null,0", func(t *testing.T) { check(t, false, null, 0) })
	t.Run("s1,0", func(t *testing.T) { check(t, true, s1, 0) })
	t.Run("s2,0", func(t *testing.T) { check(t, false, s2, 0) })
	t.Run("s3,0", func(t *testing.T) { check(t, false, s3, 0) })
	t.Run("s1,2,3,4", func(t *testing.T) { check(t, true, s1, 2, 3, 4) })
	t.Run("s2,2,3,4", func(t *testing.T) { check(t, false, s2, 2, 3, 4) })
	t.Run("s3,2,3,4", func(t *testing.T) { check(t, false, s3, 2, 3, 4) })
}

func TestAdd(t *testing.T) {
	t.Parallel()

	check := func(t *testing.T, expected Set[E], s Set[E], key E, keys ...E) {
		s = s.Copy()
		s.Add(key, keys...)
		assert.Equal(t, expected, s)
	}

	t.Run("null,0", func(t *testing.T) { check(t, New[E](0), null, 0) })
	t.Run("null,0,0", func(t *testing.T) { check(t, New[E](0), null, 0, 0) })
	t.Run("s1,0", func(t *testing.T) { check(t, s1, s1, 0) })
	t.Run("s2,0", func(t *testing.T) { check(t, New[E](0, 3, 4, 5, 6, 7), s2, 0) })
	t.Run("s1,0, 0", func(t *testing.T) { check(t, s1, s1, 0, 0) })
	t.Run("s1,2,3,4", func(t *testing.T) { check(t, s1, s1, 2, 3, 4) })
	t.Run("s2,2,3,4", func(t *testing.T) { check(t, New[E](2, 3, 4, 5, 6, 7), s2, 2, 3, 4) })
	t.Run("s3,2,3,4", func(t *testing.T) { check(t, New[E](2, 3, 4, 6, 7, 8, 9, 10), s3, 2, 3, 4) })
}

func TestDel(t *testing.T) {
	t.Parallel()

	check := func(t *testing.T, expected Set[E], s Set[E], key E, keys ...E) {
		s = s.Copy()
		s.Del(key, keys...)
		assert.Equal(t, expected, s)
	}

	t.Run("null,0", func(t *testing.T) { check(t, null, null, 0) })
	t.Run("null,0,0", func(t *testing.T) { check(t, null, null, 0, 0) })
	t.Run("s1,0", func(t *testing.T) { check(t, New[E](1, 2, 3, 4), s1, 0) })
	t.Run("s2,0", func(t *testing.T) { check(t, s2, s2, 0) })
	t.Run("s2,0, 0", func(t *testing.T) { check(t, s2, s2, 0, 0) })
	t.Run("s1,2,3,4", func(t *testing.T) { check(t, New[E](0, 1), s1, 2, 3, 4) })
	t.Run("s2,2,3,4", func(t *testing.T) { check(t, New[E](5, 6, 7), s2, 2, 3, 4) })
	t.Run("s3,2,3,4", func(t *testing.T) { check(t, s3, s3, 2, 3, 4) })
}

func TestEqual(t *testing.T) {
	t.Parallel()

	check := func(t *testing.T, expected bool, a, b Set[E], s ...Set[E]) {
		assert.Equal(t, expected, Equal(a, b, s...))
	}

	t.Run("nil,nil", func(t *testing.T) { check(t, true, snil, snil) })
	t.Run("nil,null", func(t *testing.T) { check(t, false, snil, null) })
	t.Run("nil,s1", func(t *testing.T) { check(t, false, snil, s1) })
	t.Run("nil,s2", func(t *testing.T) { check(t, false, snil, s2) })
	t.Run("nil,s3", func(t *testing.T) { check(t, false, snil, s3) })

	t.Run("null,nil", func(t *testing.T) { check(t, false, null, snil) })
	t.Run("null,null", func(t *testing.T) { check(t, true, null, null) })
	t.Run("null,s1", func(t *testing.T) { check(t, false, null, s1) })
	t.Run("null,s2", func(t *testing.T) { check(t, false, null, s2) })
	t.Run("null,s3", func(t *testing.T) { check(t, false, null, s3) })

	t.Run("s1,nil", func(t *testing.T) { check(t, false, s1, snil) })
	t.Run("s1,null", func(t *testing.T) { check(t, false, s1, null) })
	t.Run("s1,s1", func(t *testing.T) { check(t, true, s1, s1) })
	t.Run("s1,s2", func(t *testing.T) { check(t, false, s1, s2) })
	t.Run("s1,s3", func(t *testing.T) { check(t, false, s1, s3) })

	t.Run("s2,nil", func(t *testing.T) { check(t, false, s2, snil) })
	t.Run("s2,null", func(t *testing.T) { check(t, false, s2, null) })
	t.Run("s2,s1", func(t *testing.T) { check(t, false, s2, s1) })
	t.Run("s2,s2", func(t *testing.T) { check(t, true, s2, s2) })
	t.Run("s2,s3", func(t *testing.T) { check(t, false, s2, s3) })

	t.Run("s3,nil", func(t *testing.T) { check(t, false, s3, snil) })
	t.Run("s3,null", func(t *testing.T) { check(t, false, s3, null) })
	t.Run("s3,s1", func(t *testing.T) { check(t, false, s3, s1) })
	t.Run("s3,s2", func(t *testing.T) { check(t, false, s3, s2) })
	t.Run("s3,s3", func(t *testing.T) { check(t, true, s3, s3) })

	t.Run("nil,nil,nil", func(t *testing.T) { check(t, true, snil, snil, snil) })
	t.Run("null,nil,nil", func(t *testing.T) { check(t, false, null, snil, snil) })
	t.Run("nil,null,nil", func(t *testing.T) { check(t, false, snil, null, snil) })
	t.Run("nil,nil,null", func(t *testing.T) { check(t, false, snil, snil, null) })
	t.Run("null,null,nil", func(t *testing.T) { check(t, false, null, null, snil) })
	t.Run("null,nil,null", func(t *testing.T) { check(t, false, null, snil, null) })
	t.Run("nil,null,null", func(t *testing.T) { check(t, false, snil, null, null) })
	t.Run("null,null,null", func(t *testing.T) { check(t, true, null, null, null) })

	t.Run("s1,s1,s1", func(t *testing.T) { check(t, true, s1, s1, s1) })
	t.Run("null,s1,s1", func(t *testing.T) { check(t, false, null, s1, s1) })
	t.Run("s1,null,s1", func(t *testing.T) { check(t, false, s1, null, s1) })
	t.Run("s1,s1,null", func(t *testing.T) { check(t, false, s1, s1, null) })
	t.Run("null,null,s1", func(t *testing.T) { check(t, false, null, null, s1) })
	t.Run("null,s1,null", func(t *testing.T) { check(t, false, null, s1, null) })
	t.Run("s1,null,null", func(t *testing.T) { check(t, false, s1, null, null) })

	t.Run("snil,s1,s1", func(t *testing.T) { check(t, false, snil, s1, s1) })
	t.Run("s1,snil,s1", func(t *testing.T) { check(t, false, s1, snil, s1) })
	t.Run("s1,s1,snil", func(t *testing.T) { check(t, false, s1, s1, snil) })
	t.Run("snil,snil,s1", func(t *testing.T) { check(t, false, snil, snil, s1) })
	t.Run("snil,s1,snil", func(t *testing.T) { check(t, false, snil, s1, snil) })
	t.Run("s1,snil,snil", func(t *testing.T) { check(t, false, s1, snil, snil) })

	t.Run("s3,s2,s2", func(t *testing.T) { check(t, false, s3, s2, s2) })
	t.Run("s2,s3,s2", func(t *testing.T) { check(t, false, s2, s3, s2) })
	t.Run("s2,s2,s3", func(t *testing.T) { check(t, false, s2, s2, s3) })
	t.Run("s3,s3,s2", func(t *testing.T) { check(t, false, s3, s3, s2) })
	t.Run("s3,s2,s3", func(t *testing.T) { check(t, false, s3, s2, s3) })
	t.Run("s2,s3,s3", func(t *testing.T) { check(t, false, s2, s3, s3) })
}

func TestSubsetOf(t *testing.T) {
	t.Parallel()

	check := func(t *testing.T, expected, proper bool, a, b Set[E]) {
		e, p := a.IsSubsetOf(b)
		assert.Equal(t, expected, e)
		assert.Equal(t, proper, p)
	}

	t.Run("nil,nil", func(t *testing.T) { check(t, true, false, snil, snil) })
	t.Run("nil,null", func(t *testing.T) { check(t, true, false, snil, null) })
	t.Run("null,nil", func(t *testing.T) { check(t, true, false, null, snil) })
	t.Run("nil,s1", func(t *testing.T) { check(t, true, true, snil, s1) })
	t.Run("null,s1", func(t *testing.T) { check(t, true, true, null, s1) })
	t.Run("s2,nil", func(t *testing.T) { check(t, false, false, s2, snil) })
	t.Run("s2,null", func(t *testing.T) { check(t, false, false, s2, null) })
	t.Run("s1,s2", func(t *testing.T) { check(t, false, false, s1, s2) })
	t.Run("s2,s1", func(t *testing.T) { check(t, false, false, s2, s1) })
	t.Run("s3,s3", func(t *testing.T) { check(t, true, false, s3, s3) })
	t.Run("{6,7},s3", func(t *testing.T) { check(t, true, true, subset3, s3) })
	t.Run("s3,{6,7}", func(t *testing.T) { check(t, false, false, s3, subset3) })
}

func TestSupersetOf(t *testing.T) {
	t.Parallel()

	check := func(t *testing.T, expected, proper bool, a, b Set[E]) {
		e, p := a.IsSupersetOf(b)
		assert.Equal(t, expected, e)
		assert.Equal(t, proper, p)
	}

	t.Run("nil,nil", func(t *testing.T) { check(t, true, false, snil, snil) })
	t.Run("nil,null", func(t *testing.T) { check(t, true, false, snil, null) })
	t.Run("null,nil", func(t *testing.T) { check(t, true, false, null, snil) })
	t.Run("nil,s1", func(t *testing.T) { check(t, false, false, snil, s1) })
	t.Run("null,s1", func(t *testing.T) { check(t, false, false, null, s1) })
	t.Run("s2,nil", func(t *testing.T) { check(t, true, true, s2, snil) })
	t.Run("s2,null", func(t *testing.T) { check(t, true, true, s2, null) })
	t.Run("s1,s2", func(t *testing.T) { check(t, false, false, s1, s2) })
	t.Run("s2,s1", func(t *testing.T) { check(t, false, false, s2, s1) })
	t.Run("s3,s3", func(t *testing.T) { check(t, true, false, s3, s3) })
	t.Run("{6,7},s3", func(t *testing.T) { check(t, false, false, subset3, s3) })
	t.Run("s3,{6,7}", func(t *testing.T) { check(t, true, true, s3, subset3) })
}

func TestUnion(t *testing.T) {
	t.Parallel()

	union_s1_s2 := New[E](0, 1, 2, 3, 4, 5, 6, 7)
	union_s2_s3 := New[E](3, 4, 5, 6, 7, 8, 9, 10)
	union_s1_s3 := New[E](0, 1, 2, 3, 4, 6, 7, 8, 9, 10)
	union_s1_s2_s3 := New[E](0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	check := func(t *testing.T, expected Set[E], a, b Set[E], sets ...Set[E]) {
		assert.Equal(t, expected, Union(a, b, sets...))

		s := a.Copy()
		s.Update(b, sets...)
		assert.Equal(t, expected, s)
	}

	t.Run("s1,null", func(t *testing.T) { check(t, s1, s1, null) })
	t.Run("null,s1", func(t *testing.T) { check(t, s1, null, s1) })
	t.Run("s1,s1", func(t *testing.T) { check(t, s1, s1, s1) })
	t.Run("s2,s2", func(t *testing.T) { check(t, s2, s2, s2) })
	t.Run("s3,s3", func(t *testing.T) { check(t, s3, s3, s3) })
	t.Run("s1,s2", func(t *testing.T) { check(t, union_s1_s2, s1, s2) })
	t.Run("s2,s1", func(t *testing.T) { check(t, union_s1_s2, s2, s1) })
	t.Run("s2,s3", func(t *testing.T) { check(t, union_s2_s3, s2, s3) })
	t.Run("s3,s2", func(t *testing.T) { check(t, union_s2_s3, s3, s2) })
	t.Run("s1,s3", func(t *testing.T) { check(t, union_s1_s3, s1, s3) })
	t.Run("s3,s1", func(t *testing.T) { check(t, union_s1_s3, s3, s1) })
	t.Run("s1,s2,s3", func(t *testing.T) { check(t, union_s1_s2_s3, s1, s2, s3) })
	t.Run("s1,s3,s2", func(t *testing.T) { check(t, union_s1_s2_s3, s1, s3, s2) })
	t.Run("s2,s1,s3", func(t *testing.T) { check(t, union_s1_s2_s3, s2, s1, s3) })
	t.Run("s2,s3,s1", func(t *testing.T) { check(t, union_s1_s2_s3, s2, s3, s1) })
	t.Run("s3,s1,s2", func(t *testing.T) { check(t, union_s1_s2_s3, s3, s1, s2) })
	t.Run("s3,s2,s1", func(t *testing.T) { check(t, union_s1_s2_s3, s3, s2, s1) })
}

func TestIntersectionDisjoint(t *testing.T) {
	t.Parallel()

	isect_s1_s2 := New[E](3, 4)
	isect_s2_s3 := New[E](6, 7)
	isect_s1_s3 := null
	isect_s1_s2_s3 := null

	check := func(t *testing.T, expected Set[E], a, b Set[E], sets ...Set[E]) {
		assert.Equal(t, expected, Intersection(a, b, sets...))

		s := a.Copy()
		s.Intersect(b, sets...)
		assert.Equal(t, expected, s)

		disjoint := Equal(expected, null)
		f := Disjoint(a, b, sets...)
		assert.Equal(t, disjoint, f)
	}

	t.Run("s1,null", func(t *testing.T) { check(t, null, s1, null) })
	t.Run("null,s1", func(t *testing.T) { check(t, null, null, s1) })
	t.Run("s1,s1", func(t *testing.T) { check(t, s1, s1, s1) })
	t.Run("s2,s2", func(t *testing.T) { check(t, s2, s2, s2) })
	t.Run("s3,s3", func(t *testing.T) { check(t, s3, s3, s3) })
	t.Run("s1,s2", func(t *testing.T) { check(t, isect_s1_s2, s1, s2) })
	t.Run("s2,s1", func(t *testing.T) { check(t, isect_s1_s2, s2, s1) })
	t.Run("s2,s3", func(t *testing.T) { check(t, isect_s2_s3, s2, s3) })
	t.Run("s3,s2", func(t *testing.T) { check(t, isect_s2_s3, s3, s2) })
	t.Run("s1,s3", func(t *testing.T) { check(t, isect_s1_s3, s1, s3) })
	t.Run("s3,s1", func(t *testing.T) { check(t, isect_s1_s3, s3, s1) })
	t.Run("s1,s2,s3", func(t *testing.T) { check(t, isect_s1_s2_s3, s1, s2, s3) })
	t.Run("s1,s3,s2", func(t *testing.T) { check(t, isect_s1_s2_s3, s1, s3, s2) })
	t.Run("s2,s1,s3", func(t *testing.T) { check(t, isect_s1_s2_s3, s2, s1, s3) })
	t.Run("s2,s3,s1", func(t *testing.T) { check(t, isect_s1_s2_s3, s2, s3, s1) })
	t.Run("s3,s1,s2", func(t *testing.T) { check(t, isect_s1_s2_s3, s3, s1, s2) })
	t.Run("s3,s2,s1", func(t *testing.T) { check(t, isect_s1_s2_s3, s3, s2, s1) })
}

func TestDifference(t *testing.T) {
	t.Parallel()

	diff_s1_s2 := New[E](0, 1, 2)
	diff_s2_s1 := New[E](5, 6, 7)
	diff_s2_s3 := New[E](3, 4, 5)
	diff_s3_s2 := New[E](8, 9, 10)
	diff_s1_s3 := New[E](0, 1, 2, 3, 4)
	diff_s3_s1 := New[E](6, 7, 8, 9, 10)
	diff_s1_s2_s3 := New[E](0, 1, 2)
	diff_s1_s3_s2 := New[E](0, 1, 2)
	diff_s2_s1_s3 := New[E](5)
	diff_s2_s3_s1 := New[E](5)
	diff_s3_s1_s2 := New[E](8, 9, 10)
	diff_s3_s2_s1 := New[E](8, 9, 10)

	check := func(t *testing.T, expected Set[E], a, b Set[E], sets ...Set[E]) {
		assert.Equal(t, expected, Difference(a, b, sets...))

		s := a.Copy()
		s.Remove(b, sets...)
		assert.Equal(t, expected, s)
	}

	t.Run("s1,null", func(t *testing.T) { check(t, s1, s1, null) })
	t.Run("null,s1", func(t *testing.T) { check(t, null, null, s1) })
	t.Run("s1,s1", func(t *testing.T) { check(t, null, s1, s1) })
	t.Run("s2,s2", func(t *testing.T) { check(t, null, s2, s2) })
	t.Run("s3,s3", func(t *testing.T) { check(t, null, s3, s3) })
	t.Run("s1,s2", func(t *testing.T) { check(t, diff_s1_s2, s1, s2) })
	t.Run("s2,s1", func(t *testing.T) { check(t, diff_s2_s1, s2, s1) })
	t.Run("s2,s3", func(t *testing.T) { check(t, diff_s2_s3, s2, s3) })
	t.Run("s3,s2", func(t *testing.T) { check(t, diff_s3_s2, s3, s2) })
	t.Run("s1,s3", func(t *testing.T) { check(t, diff_s1_s3, s1, s3) })
	t.Run("s3,s1", func(t *testing.T) { check(t, diff_s3_s1, s3, s1) })
	t.Run("s1,s2,s3", func(t *testing.T) { check(t, diff_s1_s2_s3, s1, s2, s3) })
	t.Run("s1,s3,s2", func(t *testing.T) { check(t, diff_s1_s3_s2, s1, s3, s2) })
	t.Run("s2,s1,s3", func(t *testing.T) { check(t, diff_s2_s1_s3, s2, s1, s3) })
	t.Run("s2,s3,s1", func(t *testing.T) { check(t, diff_s2_s3_s1, s2, s3, s1) })
	t.Run("s3,s1,s2", func(t *testing.T) { check(t, diff_s3_s1_s2, s3, s1, s2) })
	t.Run("s3,s2,s1", func(t *testing.T) { check(t, diff_s3_s2_s1, s3, s2, s1) })
}

func TestSymmetricDifference(t *testing.T) {
	t.Parallel()

	symdiff_s1_s2 := New[E](0, 1, 2, 5, 6, 7)
	symdiff_s2_s3 := New[E](3, 4, 5, 8, 9, 10)
	symdiff_s1_s3 := New[E](0, 1, 2, 3, 4, 6, 7, 8, 9, 10)
	symdiff_s1_s2_s3 := New[E](0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	check := func(t *testing.T, expected Set[E], a, b Set[E], sets ...Set[E]) {
		assert.Equal(t, expected, SymmetricDifference(a, b, sets...))

		s := a.Copy()
		s.SymmetricRemove(b, sets...)
		assert.Equal(t, expected, s)
	}

	t.Run("s1,null", func(t *testing.T) { check(t, s1, s1, null) })
	t.Run("null,s1", func(t *testing.T) { check(t, s1, null, s1) })
	t.Run("s1,s1", func(t *testing.T) { check(t, null, s1, s1) })
	t.Run("s2,s2", func(t *testing.T) { check(t, null, s2, s2) })
	t.Run("s3,s3", func(t *testing.T) { check(t, null, s3, s3) })
	t.Run("s1,s2", func(t *testing.T) { check(t, symdiff_s1_s2, s1, s2) })
	t.Run("s2,s1", func(t *testing.T) { check(t, symdiff_s1_s2, s2, s1) })
	t.Run("s2,s3", func(t *testing.T) { check(t, symdiff_s2_s3, s2, s3) })
	t.Run("s3,s2", func(t *testing.T) { check(t, symdiff_s2_s3, s3, s2) })
	t.Run("s1,s3", func(t *testing.T) { check(t, symdiff_s1_s3, s1, s3) })
	t.Run("s3,s1", func(t *testing.T) { check(t, symdiff_s1_s3, s3, s1) })
	t.Run("s1,s2,s3", func(t *testing.T) { check(t, symdiff_s1_s2_s3, s1, s2, s3) })
	t.Run("s1,s3,s2", func(t *testing.T) { check(t, symdiff_s1_s2_s3, s1, s3, s2) })
	t.Run("s2,s1,s3", func(t *testing.T) { check(t, symdiff_s1_s2_s3, s2, s1, s3) })
	t.Run("s2,s3,s1", func(t *testing.T) { check(t, symdiff_s1_s2_s3, s2, s3, s1) })
	t.Run("s3,s1,s2", func(t *testing.T) { check(t, symdiff_s1_s2_s3, s3, s1, s2) })
	t.Run("s3,s2,s1", func(t *testing.T) { check(t, symdiff_s1_s2_s3, s3, s2, s1) })
}
