package set

import (
	"reflect"
	"sort"
	"testing"
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
	empty := []E{}
	one := []E{1}
	five := []E{0, 1, 2, 3, 4}

	check := func(t *testing.T, expected []E, set Set[E]) {
		keys := contents(set)
		if !reflect.DeepEqual(expected, keys) {
			t.Fatal("expected", expected, "got", keys)
		}
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

func TestEqual(t *testing.T) {
	check := func(t *testing.T, expected bool, s1, s2 Set[E]) {
		equal := Equal(s1, s2)
		if expected != equal {
			t.Fatal("expected", expected, "got", equal)
		}
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
}

func TestUnion(t *testing.T) {
	union_s1_s2 := New[E](0, 1, 2, 3, 4, 5, 6, 7)
	union_s2_s3 := New[E](3, 4, 5, 6, 7, 8, 9, 10)
	union_s1_s3 := New[E](0, 1, 2, 3, 4, 6, 7, 8, 9, 10)
	union_s1_s2_s3 := New[E](0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	check := func(t *testing.T, expected Set[E], sets ...Set[E]) {
		s := Union(sets...)
		if !Equal(expected, s) {
			t.Error("expected", contents(expected), "got", contents(s))
		}

		if len(sets) > 0 {
			s := sets[0].Copy()
			s.Union(sets[1:]...)
			if !Equal(expected, s) {
				t.Error("[in-place] expected", contents(expected), "got", contents(s))
			}
		}
	}

	t.Run("none", func(t *testing.T) { check(t, null) })
	t.Run("null", func(t *testing.T) { check(t, null, null) })
	t.Run("s1", func(t *testing.T) { check(t, s1, s1) })
	t.Run("s2", func(t *testing.T) { check(t, s2, s2) })
	t.Run("s3", func(t *testing.T) { check(t, s3, s3) })
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

func TestIntersection(t *testing.T) {
	isect_s1_s2 := New[E](3, 4)
	isect_s2_s3 := New[E](6, 7)
	isect_s1_s3 := null
	isect_s1_s2_s3 := null

	check := func(t *testing.T, expected Set[E], sets ...Set[E]) {
		s := Intersection(sets...)
		if !Equal(expected, s) {
			t.Error("expected", contents(expected), "got", contents(s))
		}

		if len(sets) > 0 {
			s := sets[0].Copy()
			s.Intersect(sets[1:]...)
			if !Equal(expected, s) {
				t.Error("[in-place] expected", contents(expected), "got", contents(s))
			}
		}
	}

	t.Run("none", func(t *testing.T) { check(t, null) })
	t.Run("null", func(t *testing.T) { check(t, null, null) })
	t.Run("s1", func(t *testing.T) { check(t, s1, s1) })
	t.Run("s2", func(t *testing.T) { check(t, s2, s2) })
	t.Run("s3", func(t *testing.T) { check(t, s3, s3) })
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

	check := func(t *testing.T, expected Set[E], sets ...Set[E]) {
		s := Difference(sets...)
		if !Equal(expected, s) {
			t.Error("expected", contents(expected), "got", contents(s))
		}

		if len(sets) > 0 {
			s := sets[0].Copy()
			s.Remove(sets[1:]...)
			if !Equal(expected, s) {
				t.Error("[in-place] expected", contents(expected), "got", contents(s))
			}
		}
	}

	t.Run("none", func(t *testing.T) { check(t, null) })
	t.Run("null", func(t *testing.T) { check(t, null, null) })
	t.Run("s1", func(t *testing.T) { check(t, s1, s1) })
	t.Run("s2", func(t *testing.T) { check(t, s2, s2) })
	t.Run("s3", func(t *testing.T) { check(t, s3, s3) })
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
	symdiff_s1_s2 := New[E](0, 1, 2, 5, 6, 7)
	symdiff_s2_s3 := New[E](3, 4, 5, 8, 9, 10)
	symdiff_s1_s3 := New[E](0, 1, 2, 3, 4, 6, 7, 8, 9, 10)
	symdiff_s1_s2_s3 := New[E](0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	check := func(t *testing.T, expected Set[E], sets ...Set[E]) {
		s := SymmetricDifference(sets...)
		if !Equal(expected, s) {
			t.Error("expected", contents(expected), "got", contents(s))
		}

		if len(sets) > 0 {
			s := sets[0].Copy()
			s.SymmetricRemove(sets[1:]...)
			if !Equal(expected, s) {
				t.Error("[in-place] expected", contents(expected), "got", contents(s))
			}
		}
	}

	t.Run("none", func(t *testing.T) { check(t, null) })
	t.Run("null", func(t *testing.T) { check(t, null, null) })
	t.Run("s1", func(t *testing.T) { check(t, null, s1) })
	t.Run("s2", func(t *testing.T) { check(t, null, s2) })
	t.Run("s3", func(t *testing.T) { check(t, null, s3) })
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
