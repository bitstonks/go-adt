package main

import "sort"

import "golang.org/x/exp/constraints"
import "github.com/bitstonks/go-adt/set"

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
	s1 := set.New[int](1, 2, 3, 4, 5, 6, 7)
	s2 := set.New[int](1, 4, 7, 8, 9, 0)

	pset(s1)
	pset(s2)
	pset(set.Union(s1, s2))
	pset(set.Union(s2, s1))
	pset(set.Intersection(s2, s1))
	pset(set.Intersection(s1, s2))
	pset(set.Difference(s1, s2))
	pset(set.Difference(s2, s1))
	pset(set.SymmetricDifference(s1, s2))
	pset(set.SymmetricDifference(s2, s1))
}
