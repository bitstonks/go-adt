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

// The union of all the sets: ⋃(sets) = set[0] ∪ set[1] ∪ set[2] ...
func Union[Key comparable](sets ...Set[Key]) Set[Key] {
	resultset := make(Set[Key])
	resultset.Union(sets...)
	return resultset
}

// Like Union, but modifies the set in place.
func (s Set[Key]) Union(sets ...Set[Key]) {
	if sets == nil {
		return
	}
	for i := range sets {
		for k := range sets[i] {
			s[k] = struct{}{}
		}
	}
}

// The intersection of all the sets: ⋂(sets) = set[0] ∩ set[1] ∩ set[2] ...
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

// Like Intersection, but modifies the set in place.
func (s Set[Key]) Intersect(sets ...Set[Key]) {
	if sets == nil {
		return
	}

	// Clear the set if any of the arguments is an empty set.
	for i := range sets {
		if len(sets[i]) == 0 {
			for k := range s {
				delete(s, k)
			}
			return
		}
	}

	rm := make([]Key, 0, len(s))
outer:
	for k := range s {
		for i := range sets {
			if !sets[i].Contains(k) {
				rm = append(rm, k)
				continue outer
			}
		}
	}
	for i := range rm {
		delete(s, rm[i])
	}
}

// The difference of all the sets: sets[0] ∖ sets[1] ∖ sets[2] ...
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

// Like Difference, but modifies the set in place.
func (s Set[Key]) Remove(sets ...Set[Key]) {
	if sets == nil {
		return
	}

	rm := make([]Key, 0, len(s))
outer:
	for k := range s {
		for i := range sets {
			if sets[i].Contains(k) {
				rm = append(rm, k)
				continue outer
			}
		}
	}
	for i := range rm {
		delete(s, rm[i])
	}
}

// Symmetric difference of all the sets: ⋃(sets) ∖ ⋂(sets).
func SymmetricDifference[Key comparable](sets ...Set[Key]) Set[Key] {
	// The symmetric difference of no sets is the empty set
	if sets == nil {
		return make(Set[Key])
	}
	return Difference(Union(sets...), Intersection(sets...))
}

// Like SymmetricDifference, but modifies the set in place.
func (s Set[Key]) SymmetricRemove(sets ...Set[Key]) {
	if sets == nil {
		return
	}
	rm := Intersection(sets...)
	rm.Intersect(s)
	s.Union(sets...)
	s.Remove(rm)
}
