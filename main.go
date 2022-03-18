package main

import (
	"github.com/bitstonks/go-adt/set"
	"golang.org/x/exp/constraints"
	"sort"
)

func pset[K constraints.Ordered](s set.Set[K]) {
	i := 0
	keys := make([]K, len(s))
	for e := range s {
		keys[i] = e
		i += 1
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	sep := ""
	for i := range keys {
		print(sep, keys[i])
		sep = " "
	}
	print("\n")
}

func main() {
	var s set.Set[int]
	s1 := set.New[int](1, 2, 3, 4, 5, 6, 7)
	s2 := set.New[int](1, 0, 0, 4, 0, 0, 7, 8, 9)

	print("s1 = ")
	pset(s1)
	print("s2 = ")
	pset(s2)

	print("\nUnion: s1 ∪ s2\n")
	pset(set.Union(s1, s2))
	s = s1.Copy()
	s.Union(s2)
	pset(s)

	print("\nUnion: s2 ∪ s1\n")
	pset(set.Union(s2, s1))
	s = s2.Copy()
	s.Union(s1)
	pset(s)

	print("\nIntersection s1 ∩ s2:\n")
	pset(set.Intersection(s1, s2))
	s = s1.Copy()
	s.Intersect(s2)
	pset(s)

	print("\nIntersection s2 ∩ s1:\n")
	pset(set.Intersection(s2, s1))
	s = s2.Copy()
	s.Intersect(s1)
	pset(s)

	print("\nDifference: s1 ∖ s2\n")
	pset(set.Difference(s1, s2))
	s = s1.Copy()
	s.Remove(s2)
	pset(s)

	print("\nDifference: s2 ∖ s1\n")
	pset(set.Difference(s2, s1))
	s = s2.Copy()
	s.Remove(s1)
	pset(s)

	print("\nSymmetric difference: s1 ⊖ s2\n")
	pset(set.SymmetricDifference(s1, s2))
	s = s1.Copy()
	s.SymmetricRemove(s2)
	pset(s)

	print("\nSymmetric difference: s2 ⊖ s1\n")
	pset(set.SymmetricDifference(s2, s1))
	s = s2.Copy()
	s.SymmetricRemove(s1)
	pset(s)
}
