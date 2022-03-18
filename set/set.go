package set

import "sort"

type Set[Key comparable] map[Key]struct{}

func (s Set[Key]) Add(keys ...Key) {
	if keys == nil {
		return
	}
	for i := range keys {
		s[keys[i]] = struct{}{}
	}
}

func (s Set[Key]) Contains(key Key) bool {
	_, exists := s[key]
	return exists
}

func (s Set[Key]) Copy() Set[Key] {
	resultset := make(Set[Key], len(s))
	for k := range s {
		resultset[k] = struct{}{}
	}
	return resultset
}

func New[Key comparable](keys ...Key) Set[Key] {
	if keys == nil {
		return make(Set[Key])
	}
	resultset := make(Set[Key], len(keys))
	resultset.Add(keys...)
	return resultset
}

func Union[Key comparable](sets ...Set[Key]) Set[Key] {
	resultset := make(Set[Key])
	for i := range sets {
		for k := range sets[i] {
			resultset[k] = struct{}{}
		}
	}
	return resultset
}

func Intersection[Key comparable](sets ...Set[Key]) Set[Key] {
	// The intersection of no sets is the empty set.
	if sets == nil {
		return make(Set[Key])
	}

	// The intersection of one set is the set itself.
	if len(sets) == 1 {
		return sets[0].Copy()
	}

	// Use the smallest set as the candidate result.
	// Do *not* modify the function argument array, copy it before sorting.
	sorted_sets := append(make([]Set[Key], 0, len(sets)), sets...)
	sort.Slice(sorted_sets, func(i, j int) bool {
		return len(sorted_sets[i]) < len(sorted_sets[j])
	})
	candidate := sorted_sets[0]

	// Any empty set in the arguments produces an empty intersection.
	if len(candidate) == 0 {
		return make(Set[Key])
	}

	others := sorted_sets[1:]
	resultset := make(Set[Key], len(candidate))
outer:
	for k := range candidate {
		for i := range others {
			if !others[i].Contains(k) {
				continue outer
			}
		}
		resultset[k] = struct{}{}
	}
	return resultset
}

func Difference[Key comparable](sets ...Set[Key]) Set[Key] {
	// The difference of no sets is the empty set
	if sets == nil {
		return make(Set[Key])
	}

	// The difference of one set is the set itself.
	if len(sets) == 1 {
		return sets[0].Copy()
	}

	candidate := sets[0]
	others := sets[1:]
	resultset := make(Set[Key], len(candidate))
outer:
	for k := range candidate {
		for i := range others {
			if others[i].Contains(k) {
				continue outer
			}
		}
		resultset[k] = struct{}{}
	}
	return resultset
}

func SymmetricDifference[Key comparable](sets ...Set[Key]) Set[Key] {
	return Difference(Union(sets...), Intersection(sets...))
}
