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

// Checks if sets are equal: ⋂(a, b, sets...) = ⋃(a, b, sets...)
func Equal[Key comparable](a, b Set[Key], sets ...Set[Key]) bool {
	// The same set is always equal to itself.
	if len(sets) == 0 && sameobject(a, b) {
		return true
	}

	if !equal(a, b) {
		return false
	}
	for i := range sets {
		if !equal(a, sets[i]) {
			return false
		}
	}
	return true
}

// Checks if sets are disjoint: ⋂(a, b, sets) = ∅
func Disjoint[Key comparable](a, b Set[Key], sets ...Set[Key]) bool {
	// The same set is never disjoint against itself unless it's the empty set.
	if len(sets) == 0 && (sameobject(a, b) || (len(a) == 0 && len(b) == 0)) {
		return (len(a) == 0)
	}

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

// Returns a pair of values:
//  - ok=true if s is a subset of other
//  - proper=true if s is a proper subset (i.e., !Equal(s, other)))
func (s Set[Key]) IsSubsetOf(other Set[Key]) (ok bool, proper bool) {
	// The same set is always its own improper subset.
	if sameobject(s, other) {
		return true, false
	}

	if len(s) > len(other) {
		return false, false
	}
	if len(s) > 0 {
		for k := range s {
			if !other.has(k) {
				return false, false
			}
		}
	}
	return true, len(s) < len(other)
}

// Returns a pair of values:
//  - ok=true if s is a superset of other
//  - proper=true if s is a proper superset (i.e., !Equal(s, other)))
func (s Set[Key]) IsSupersetOf(other Set[Key]) (ok bool, proper bool) {
	return other.IsSubsetOf(s)
}

// Creates a deep copy of the set. Will never return nil.
func (s Set[Key]) Copy() Set[Key] {
	resultset := make(Set[Key], len(s))
	for k := range s {
		resultset[k] = struct{}{}
	}
	return resultset
}

// Checks if the set contains all of the given keys.
func (s Set[Key]) Contains(key Key, keys ...Key) bool {
	// An empty set contains no keys.
	if len(s) == 0 {
		return false
	}

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

// Delets keys from the set.
func (s Set[Key]) Del(key Key, keys ...Key) {
	// An empty set contains no keys.
	if len(s) == 0 {
		return
	}

	delete(s, key)
	for i := range keys {
		if len(s) == 0 {
			return
		}
		delete(s, keys[i])
	}
}

// The union of all the sets: ⋃(a, b, sets) = a ∪ b ∪ sets[0] ∪ sets[1] ...
func Union[Key comparable](a, b Set[Key], sets ...Set[Key]) Set[Key] {
	// The union of a set with itself is the set.
	if len(sets) == 0 && sameobject(a, b) {
		return a.Copy()
	}

	resultset := a.Copy()
	resultset.Update(b, sets...)
	return resultset
}

// Like Union, but modifies the set in place.
func (s Set[Key]) Update(a Set[Key], sets ...Set[Key]) {
	// The union of a set with itself is the set.
	if len(sets) == 0 && sameobject(s, a) {
		return
	}

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
	// The result will be empty if any of the sets are empty.
	if len(a) == 0 || len(b) == 0 {
		return make(Set[Key])
	}

	// The intersection of a set with itself is the set.
	if len(sets) == 0 && sameobject(a, b) {
		return a.Copy()
	}

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

	// The intersection of a set with itself is the set.
	if len(sets) == 0 && sameobject(s, a) {
		return
	}

	// The intersection of a set with an empty set is empty.
	if len(a) == 0 {
		s.clear()
		return
	}

	// Clear the set if any of the arguments is an empty set.
	for i := range sets {
		if len(sets[i]) == 0 {
			s.clear()
			return
		}
	}

	sets = append(sets, a)
	rm := make([]Key, 0, len(s))
outer:
	for k := range s {
		for i := range sets {
			if !sets[i].has(k) {
				rm = append(rm, k)
				if len(rm) == len(s) {
					// The result will be empty.
					break outer
				}
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
	// The result will be empty if the first set is empty or when we're finding
	// the difference of the same set.
	if len(a) == 0 || sameobject(a, b) {
		return make(Set[Key])
	}

	// The difference with a null set is the same set.
	if len(b) == 0 && len(sets) == 0 {
		return a.Copy()
	}

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
	// The result will be empty if this set is empty
	if len(s) == 0 {
		return
	}

	// The result will be empty if we're removing the set from itself.
	if sameobject(s, a) {
		s.clear()
		return
	}

	// The result will be unchanged if the other sets are empty.
	if len(a) == 0 && len(sets) == 0 {
		return
	}

	sets = append(sets, a)
	rm := make([]Key, 0, len(s))
outer:
	for k := range s {
		for i := range sets {
			if sets[i].has(k) {
				rm = append(rm, k)
				if len(rm) == len(s) {
					// The result will be empty.
					break outer
				}
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
	// The symmetric difference of a set with itself is the empty set.
	if len(sets) == 0 && sameobject(a, b) {
		return make(Set[Key])
	}

	return Difference(Union(a, b, sets...), Intersection(a, b, sets...))
}

// Like SymmetricDifference, but modifies the set in place.
func (s Set[Key]) SymmetricRemove(a Set[Key], sets ...Set[Key]) {
	// The symmetric difference of a set with itself is the empty set.
	if len(sets) == 0 && sameobject(s, a) {
		s.clear()
		return
	}

	rm := Intersection(s, a, sets...)
	s.Update(a, sets...)
	s.Remove(rm)
}

// Does the same check as reflect.DeepEqual() for maps.
// See: https://github.com/golang/go/blob/master/src/reflect/deepequal.go
func sameobject[Key comparable](a, b Set[Key]) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	return va.UnsafePointer() == vb.UnsafePointer()
}

// All sets of zero length are equal, regardless of representation.
func equal[Key comparable](a, b Set[Key]) bool {
	return (len(a) == 0 && len(b) == 0) || reflect.DeepEqual(a, b)
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
