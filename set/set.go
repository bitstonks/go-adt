package set

import (
	"reflect"
	"sort"
)

type Set[Key comparable] map[Key]struct{}

// Creates a new set that contains all the given keys.
func New[Key comparable](keys ...Key) Set[Key] {
	if len(keys) == 0 {
		return make(Set[Key])
	}
	resultset := make(Set[Key], len(keys))
	for i := range keys {
		resultset[keys[i]] = struct{}{}
	}
	return resultset
}

// Checks if sets are equal: ⋂(a, b, sets) = a
func Equal[Key comparable](a, b Set[Key], sets ...Set[Key]) bool {
	if !reflect.DeepEqual(a, b) {
		return false
	}
	for i := range sets {
		if !reflect.DeepEqual(a, sets[i]) {
			return false
		}
	}
	return true
}

// Checks if sets are disjoint: ⋂(a, b, sets) = ∅
func Disjoint[Key comparable](a, b Set[Key], sets ...Set[Key]) bool {
	// Use the smallest set to check the others.
	sorted_sets := sortSetsByLength(a, b, sets...)
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

// Creates a deep copy of the set.
func (s Set[Key]) Copy() Set[Key] {
	resultset := make(Set[Key], len(s))
	for k := range s {
		resultset[k] = struct{}{}
	}
	return resultset
}

// Checks if the set contains all of the given keys.
func (s Set[Key]) Contains(key Key, keys ...Key) bool {
	if !s.has(key) {
		return false
	}
	for i := range keys {
		if !s.has(keys[i]) {
			return false
		}
	}
	return true
}

// Adds keys to the set.
func (s Set[Key]) Add(key Key, keys ...Key) {
	s[key] = struct{}{}
	for i := range keys {
		s[keys[i]] = struct{}{}
	}
}

// The union of all the sets: ⋃(a, b, sets) = a ∪ b ∪ sets[0] ∪ sets[1] ...
func Union[Key comparable](a, b Set[Key], sets ...Set[Key]) Set[Key] {
	resultset := a.Copy()
	resultset.Extend(b, sets...)
	return resultset
}

// Like Union, but modifies the set in place.
func (s Set[Key]) Extend(a Set[Key], sets ...Set[Key]) {
	for k := range a {
		s[k] = struct{}{}
	}
	for i := range sets {
		for k := range sets[i] {
			s[k] = struct{}{}
		}
	}
}

// The intersection of all the sets: ⋂(a, b, sets) = a ∩ b ∩ sets[0] ∩ sets[1] ...
func Intersection[Key comparable](a, b Set[Key], sets ...Set[Key]) Set[Key] {
	// Use the smallest set as the candidate result.
	sorted_sets := sortSetsByLength(a, b, sets...)
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
func (s Set[Key]) Intersect(a Set[Key], sets ...Set[Key]) {
	// The result will be empty if this set is empty.
	if len(s) == 0 {
		return
	}

	// Clear the set if any of the arguments is an empty set.
	sets = append(sets, a)
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

// The difference of all the sets: a ∖ b ∖ sets[0] ∖ sets[1] ...
func Difference[Key comparable](a, b Set[Key], sets ...Set[Key]) Set[Key] {
	sets = append(sets, b)
	resultset := make(Set[Key], len(a))
outer:
	for k := range a {
		for i := range sets {
			if sets[i].has(k) {
				continue outer
			}
		}
		resultset[k] = struct{}{}
	}
	return resultset
}

// Like Difference, but modifies the set in place.
func (s Set[Key]) Remove(a Set[Key], sets ...Set[Key]) {
	sets = append(sets, a)
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

// Symmetric difference of all the sets: ⋃(a, b, sets) ∖ ⋂(a, b, sets).
func SymmetricDifference[Key comparable](a, b Set[Key], sets ...Set[Key]) Set[Key] {
	return Difference(Union(a, b, sets...), Intersection(a, b, sets...))
}

// Like SymmetricDifference, but modifies the set in place.
func (s Set[Key]) SymmetricRemove(a Set[Key], sets ...Set[Key]) {
	// The symmetric difference of a set with itself is the empty set.
	if &s == &a {
		s.clear()
		return
	}

	rm := s.Copy()
	rm.Intersect(a, sets...)
	s.Extend(a, sets...)
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

func sortSetsByLength[Key comparable](a, b Set[Key], sets ...Set[Key]) []Set[Key] {
	// Do *not* modify the function argument array, copy it before sorting.
	if len(sets) == 0 {
		if len(b) < len(a) {
			return []Set[Key]{b, a}
		}
		return []Set[Key]{a, b}
	} else {
		sorted := make([]Set[Key], 0, 2+len(sets))
		sorted = append(append(append(sorted, a), b), sets...)
		sort.Slice(sorted, func(i, j int) bool {
			return len(sorted[i]) < len(sorted[j])
		})
		return sorted
	}
}
