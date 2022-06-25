package main

import (
	"github.com/bitstonks/go-adt/set"

	"sort"
)

type ynt struct{ n int }

func pset(s set.Set[ynt]) {
	i := 0
	keys := make([]ynt, len(s))
	for e := range s {
		keys[i] = e
		i += 1
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].n < keys[j].n })

	sep := ""
	for i := range keys {
		print(sep, keys[i].n)
		sep = " "
	}
	print("\n")
}

func main() {
	var s set.Set[ynt]
	s1 := set.New[ynt](ynt{1}, ynt{2}, ynt{3}, ynt{4}, ynt{5}, ynt{6}, ynt{7})
	s2 := set.New[ynt](ynt{1}, ynt{0}, ynt{0}, ynt{4}, ynt{0}, ynt{0}, ynt{7}, ynt{8}, ynt{9})

	print("s1 = ")
	pset(s1)
	print("s2 = ")
	pset(s2)

	print("\nUnion: s1 ∪ s2\n")
	pset(set.Union(s1, s2))
	s = s1.Copy()
	s.Update(s2)
	pset(s)

	print("\nUnion: s2 ∪ s1\n")
	pset(set.Union(s2, s1))
	s = s2.Copy()
	s.Update(s1)
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
