# go-adt

Go implementations of different abstract data types using generics.
Requires Go 1.18+.

 * `./set`: generic set

## set

Implements a generic set whose keys must fulfill the `comparable` constraint.
The implementation uses Go's `map` with an empty `struct{}` element.

Construction:
  * `set.New[T]()` creates an empty set; equivalent to `make(map[T]struct{})`.
  * `set.New(keys...)` creates a new set filled with _keys_. Equivalent to
    `set.New[T]().Add(keys...)`.
  * `s.Copy()` creates a new set that is equal to _s_.

Queries and comparison:
  * `len(s)` returns the length of the set.
  * `s.Contains(keys...)` returns `true` if _s_ contains all the _keys_.
    Equivalent to repeating `_, ok := s[key]` for each key.
  * `set.Equal(sets...)` returns `true` if all the _sets_ are of equal size and
    contain the same keys: ⋂(_sets_) = ⋃(_sets_). The `nil` and empty sets
    are considered equal.
  * `set.Disjoint(sets...)` returns `true` if all the _sets_ have no keys in
    common: ⋂(_sets_) = ∅, unless all _sets_ are empty.
  * `s1.isSubset(s2)` checks if _s_<sub>1</sub> is a (proper) subset of
    _s_<sub>2</sub>: _s_<sub>1</sub> ⊆ _s_<sub>2</sub>.
  * `s1.isSuperset(s2)` is the opposite of `.isSubset()`: _s_<sub>1</sub> ⊇
    _s_<sub>2</sub>.

Non-modifying operations:
  * `set.Union(sets...)` creates a new set that is the union of all the _sets_.
  * `set.Intersection(sets...)` creates a new set that is the intersection of
    all the _sets_.
  * `set.Difference(sets...)` creates a new set that is the sequential
    difference of all the _sets_: _s_<sub>1</sub> ∖ _s_<sub>2</sub> ∖
    _s_<sub>3</sub>...
  * `set.SymmetricDifference(sets...)` creates a new set that is the difference
    between the union and intersection of the _sets_: _s_<sub>1</sub> ⊖
    _s_<sub>2</sub> ⊖ _s_<sub>3</sub>... = ⋃(_sets_) ∖ ⋂(_sets_).

Modifying operations:
  * `s.Add(keys...)` updates _s_ in place by adding all _keys_ to it.
  * `s.Del(keys...)` updates _s_ in place by removing all _keys_ from it.
  * `s.Update(sets...)` is equivalent to `s = set.Union(s, sets...)`.
  * `s.Intersect(sets...)` is equivalent to `s = set.Intersection(s, sets...)`.
  * `s.Remove(sets...)` is equivalent to `s = set.Difference(s, sets...)`.
  * `s.SymmetricRemove(sets...)` is equivalent to `s = set.SymmetricDifference(s, sets...)`.
