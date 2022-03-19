package set

import (
	"reflect"
	"sort"
)

type Set[Key comparable] map[Key]struct{}

func (s Set[Key]) Add(keys ...Key) {
	if len(keys) == 0 {
		return
	}
	for i := range keys {
		s[keys[i]] = struct{}{}
	}
}

// Checks if the set contains all of the given keys.
// Returns true when called without arguments.
func (s Set[Key]) Contains(keys ...Key) bool {
	for i := range keys {
		if !s.has(keys[i]) {
			return false
		}
	}
	return true
}

func (s Set[Key]) Copy() Set[Key] {
	resultset := make(Set[Key], len(s))
	for k := range s {
		resultset[k] = struct{}{}
	}
	return resultset
}

func New[Key comparable](keys ...Key) Set[Key] {
	if len(keys) == 0 {
		return make(Set[Key])
	}
	resultset := make(Set[Key], len(keys))
	resultset.Add(keys...)
	return resultset
}

// Checks if sets are equal: ⋃(sets) = sets[0]
// No sets and one set are always equal.
func Equal[Key comparable](sets ...Set[Key]) bool {
	if len(sets) > 1 {
		first := sets[0]
		rest := sets[1:]
		for i := range rest {
			if !reflect.DeepEqual(first, rest[i]) {
				return false
			}
		}
	}
	return true
}

// Checks if sets are disjoint: ⋂(sets) = ∅
func Disjoint[Key comparable](sets ...Set[Key]) bool {
	// No sets are always disjoint.
	if len(sets) == 0 {
		return true
	}

	// One set is disjoint if it's empty.
	if len(sets) == 1 {
		return len(sets[0]) == 0
	}

	// Use the smallest set to check the others.
	sorted_sets := sortSetsByLength(sets...)
	candidate := sorted_sets[0]

	// Any empty set in the arguments means the sets are disjoint.
	if len(candidate) == 0 {
		return true
	}

	others := sorted_sets[1:]
outer:
	for k := range candidate {
		for i := range others {
			if !others[i].has(k) {
				continue outer
			}
		}
		return false
	}
	return true
}

// The union of all the sets: ⋃(sets) = sets[0] ∪ sets[1] ∪ sets[2] ...
func Union[Key comparable](sets ...Set[Key]) Set[Key] {
	resultset := make(Set[Key])
	resultset.Union(sets...)
	return resultset
}

// Like Union, but modifies the set in place.
func (s Set[Key]) Union(sets ...Set[Key]) {
	if len(sets) == 0 {
		return
	}
	for i := range sets {
		for k := range sets[i] {
			s[k] = struct{}{}
		}
	}
}

// The intersection of all the sets: ⋂(sets) = sets[0] ∩ sets[1] ∩ sets[2] ...
func Intersection[Key comparable](sets ...Set[Key]) Set[Key] {
	// The intersection of no sets is the empty set.
	if len(sets) == 0 {
		return make(Set[Key])
	}

	// The intersection of one set is the set itself.
	if len(sets) == 1 {
		return sets[0].Copy()
	}

	// Use the smallest set as the candidate result.
	sorted_sets := sortSetsByLength(sets...)
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
			if !others[i].has(k) {
				continue outer
			}
		}
		resultset[k] = struct{}{}
	}
	return resultset
}

// Like Intersection, but modifies the set in place.
func (s Set[Key]) Intersect(sets ...Set[Key]) {
	if len(sets) == 0 {
		return
	}

	// Clear the set if any of the arguments is an empty set.
	for i := range sets {
		if len(sets[i]) == 0 {
			s.clear()
			return
		}
	}

	rm := make([]Key, 0, len(s))
outer:
	for k := range s {
		for i := range sets {
			if !sets[i].has(k) {
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
	// The difference of no sets is the empty set.
	if len(sets) == 0 {
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
			if others[i].has(k) {
				continue outer
			}
		}
		resultset[k] = struct{}{}
	}
	return resultset
}

// Like Difference, but modifies the set in place.
func (s Set[Key]) Remove(sets ...Set[Key]) {
	// The difference of one set is the set itself.
	if len(sets) == 0 {
		return
	}

	rm := make([]Key, 0, len(s))
outer:
	for k := range s {
		for i := range sets {
			if sets[i].has(k) {
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
	// The symmetric difference of no sets or one set is the empty set.
	if len(sets) == 0 || len(sets) == 1 {
		return make(Set[Key])
	}

	return Difference(Union(sets...), Intersection(sets...))
}

// Like SymmetricDifference, but modifies the set in place.
func (s Set[Key]) SymmetricRemove(sets ...Set[Key]) {
	// The symmetric difference of a set with itself is the empty set.
	if len(sets) == 0 {
		s.clear()
		return
	}

	rm := s.Copy()
	rm.Intersect(sets...)
	s.Union(sets...)
	s.Remove(rm)
}

func (s Set[Key]) clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Set[Key]) has(key Key) bool {
	_, exists := s[key]
	return exists
}

func sortSetsByLength[Key comparable](sets ...Set[Key]) []Set[Key] {
	if len(sets) < 2 {
		panic("invalid arguments")
	}

	// Do *not* modify the function argument array, copy it before sorting.
	if len(sets) == 2 {
		if len(sets[1]) < len(sets[0]) {
			return []Set[Key]{sets[1], sets[0]}
		}
		return sets
	} else {
		sorted := append(make([]Set[Key], 0, len(sets)), sets...)
		sort.Slice(sorted, func(i, j int) bool {
			return len(sorted[i]) < len(sorted[j])
		})
		return sorted
	}
}
