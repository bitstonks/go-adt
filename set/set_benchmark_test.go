package set

import (
	"math/rand"
	"testing"
	"time"
)

type tkey = int64
type tset = Set[tkey]

func randomKeyArray(size uint) []tkey {
	s := rand.NewSource(time.Now().UnixMicro())
	a := make([]tkey, 0, size)
	for i := uint(0); i < size; i++ {
		a = append(a, s.Int63())
	}
	return a
}

var (
	randomKeys_1000 []tkey = randomKeyArray(1000)
	nullSet         tset   = New[tkey]()
	randomSetA_1000 tset   = New(randomKeys_1000...)
	randomSetB_1000 tset   = New(randomKeys_1000...)
	randomSetX_1000 tset   = New(randomKeyArray(1000)...)
	randomSetX_2000 tset   = New(randomKeyArray(2000)...)
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[tkey]()
	}
}

func BenchmarkNewFill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(randomKeys_1000...)
	}
}

func BenchmarkEqual_NullNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Equal(nullSet, nil)
	}
}

func BenchmarkEqual_WithNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Equal(randomSetX_1000, nullSet)
	}
}

func BenchmarkEqual_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Equal(nullSet, randomSetX_1000)
	}
}

func BenchmarkEqual_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Equal(randomSetA_1000, randomSetA_1000)
	}
}

func BenchmarkEqual_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Equal(randomSetA_1000, randomSetB_1000)
	}
}

func BenchmarkEqual_SameSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Equal(randomSetA_1000, randomSetX_1000)
	}
}

func BenchmarkEqual_DifferentSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Equal(randomSetA_1000, randomSetX_2000)
	}
}

func BenchmarkDisjoint_NullNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Disjoint(nullSet, nil)
	}
}

func BenchmarkDisjoint_WithNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Disjoint(randomSetA_1000, nullSet)
	}
}

func BenchmarkDisjoint_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Disjoint(nullSet, randomSetA_1000)
	}
}

func BenchmarkDisjoint_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Disjoint(randomSetA_1000, randomSetA_1000)
	}
}

func BenchmarkDisjoint_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Disjoint(randomSetA_1000, randomSetB_1000)
	}
}

func BenchmarkDisjoint_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Disjoint(randomSetA_1000, randomSetX_1000)
	}
}

func BenchmarkIsSubsetOf_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nullSet.IsSubsetOf(randomSetA_1000)
	}
}

func BenchmarkIsSubsetOf_WithNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetA_1000.IsSubsetOf(nullSet)
	}
}

func BenchmarkIsSubsetOf_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetA_1000.IsSubsetOf(randomSetA_1000)
	}
}

func BenchmarkIsSubsetOf_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetA_1000.IsSubsetOf(randomSetB_1000)
	}
}

func BenchmarkIsSubsetOf_SameSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetA_1000.IsSubsetOf(randomSetX_1000)
	}
}

func BenchmarkIsSubsetOf_SmallerSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetX_2000.IsSubsetOf(randomSetX_1000)
	}
}

func BenchmarkIsSubsetOf_LargerSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetX_1000.IsSubsetOf(randomSetX_2000)
	}
}

func BenchmarkCopy_Null(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nullSet.Copy()
	}
}

func BenchmarkCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetX_1000.Copy()
	}
}

func BenchmarkElements_Null(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nullSet.Elements()
	}
}

func BenchmarkElements(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetX_1000.Elements()
	}
}

func BenchmarkContains_Null(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nullSet.Contains(0, randomKeys_1000...)
	}
}

func BenchmarkContains_SameKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetA_1000.Contains(0, randomKeys_1000...)
	}
}

func BenchmarkContains_DifferentKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSetX_1000.Contains(0, randomKeys_1000...)
	}
}

func BenchmarkAdd_Null(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[tkey]().Add(0, randomKeys_1000...)
	}
}

func BenchmarkAdd_SameKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.Add(0, randomKeys_1000...)
	}
}

func BenchmarkAdd_DifferentKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetX_1000.Copy()
		b.StartTimer()
		s.Add(0, randomKeys_1000...)
	}
}

func BenchmarkDel_Null(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[tkey]().Del(0, randomKeys_1000...)
	}
}

func BenchmarkDel_SameKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.Del(0, randomKeys_1000...)
	}
}

func BenchmarkDel_DifferentKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetX_1000.Copy()
		b.StartTimer()
		s.Del(0, randomKeys_1000...)
	}
}

func BenchmarkUnion_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Union(nullSet, randomSetX_1000)
	}
}

func BenchmarkUnion_WithNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Union(randomSetX_1000, nullSet)
	}
}

func BenchmarkUnion_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Union(randomSetA_1000, randomSetA_1000)
	}
}

func BenchmarkUnion_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Union(randomSetA_1000, randomSetB_1000)
	}
}

func BenchmarkUnion_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Union(randomSetA_1000, randomSetX_1000)
	}
}

func BenchmarkUpdate_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[tkey]().Update(randomSetX_1000)
	}
}

func BenchmarkUpdate_WithNull(b *testing.B) {
	s := randomSetA_1000.Copy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Update(nullSet)
	}
}

func BenchmarkUpdate_SameObject(b *testing.B) {
	s := randomSetA_1000.Copy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Update(s)
	}
}

func BenchmarkUpdate_SameSet(b *testing.B) {
	s := randomSetA_1000.Copy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Update(randomSetA_1000)
	}
}

func BenchmarkUpdate_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.Update(randomSetX_1000)
	}
}

func BenchmarkIntersection_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Intersection(nullSet, randomSetX_1000)
	}
}

func BenchmarkIntersection_WithNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Intersection(randomSetX_1000, nullSet)
	}
}

func BenchmarkIntersection_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Intersection(randomSetA_1000, randomSetA_1000)
	}
}

func BenchmarkIntersection_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Intersection(randomSetA_1000, randomSetB_1000)
	}
}

func BenchmarkIntersection_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Intersection(randomSetA_1000, randomSetX_1000)
	}
}

func BenchmarkIntersect_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[tkey]().Intersect(randomSetX_1000)
	}
}

func BenchmarkIntersect_WithNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetX_1000.Copy()
		b.StartTimer()
		s.Intersect(nullSet)
	}
}

func BenchmarkIntersect_SameObject(b *testing.B) {
	s := randomSetA_1000.Copy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Intersect(s)
	}
}

func BenchmarkIntersect_SameSet(b *testing.B) {
	s := randomSetA_1000.Copy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Intersect(randomSetA_1000)
	}
}

func BenchmarkIntersect_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.Intersect(randomSetX_1000)
	}
}

func BenchmarkDifference_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Difference(nullSet, randomSetX_1000)
	}
}

func BenchmarkDifference_WithNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Difference(randomSetX_1000, nullSet)
	}
}

func BenchmarkDifference_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Difference(randomSetA_1000, randomSetA_1000)
	}
}

func BenchmarkDifference_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Difference(randomSetA_1000, randomSetB_1000)
	}
}

func BenchmarkDifference_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Difference(randomSetA_1000, randomSetX_1000)
	}
}

func BenchmarkRemove_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[tkey]().Remove(randomSetX_1000)
	}
}

func BenchmarkRemove_WithNull(b *testing.B) {
	s := randomSetA_1000.Copy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Remove(nullSet)
	}
}

func BenchmarkRemove_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.Remove(s)
	}
}

func BenchmarkRemove_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.Remove(randomSetA_1000)
	}
}

func BenchmarkRemove_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.Remove(randomSetX_1000)
	}
}

func BenchmarkSymmetricDifference_NullNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SymmetricDifference(nullSet, nil)
	}
}

func BenchmarkSymmetricDifference_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SymmetricDifference(nullSet, randomSetX_1000)
	}
}

func BenchmarkSymmetricDifference_WithNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SymmetricDifference(randomSetX_1000, nullSet)
	}
}

func BenchmarkSymmetricDifference_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SymmetricDifference(randomSetA_1000, randomSetA_1000)
	}
}

func BenchmarkSymmetricDifference_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SymmetricDifference(randomSetA_1000, randomSetB_1000)
	}
}

func BenchmarkSymmetricDifference_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SymmetricDifference(randomSetA_1000, randomSetX_1000)
	}
}

func BenchmarkSymmetricRemove_NullNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[tkey]().SymmetricRemove(nil)
	}
}

func BenchmarkSymmetricRemove_NullWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[tkey]().SymmetricRemove(randomSetX_1000)
	}
}

func BenchmarkSymmetricRemove_WithNull(b *testing.B) {
	s := randomSetA_1000.Copy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SymmetricRemove(nullSet)
	}
}

func BenchmarkSymmetricRemove_SameObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.SymmetricRemove(s)
	}
}

func BenchmarkSymmetricRemove_SameSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.SymmetricRemove(randomSetA_1000)
	}
}

func BenchmarkSymmetricRemove_DifferentSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := randomSetA_1000.Copy()
		b.StartTimer()
		s.SymmetricRemove(randomSetX_1000)
	}
}
