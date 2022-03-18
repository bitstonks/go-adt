package adt

var keyExists = struct{}{}

type Set[E comparable] map[E]struct{}

func NewSet[E comparable](xs ...E) Set[E] {
	s := make(Set[E])
	for _, x := range xs {
		s[x] = keyExists
	}
	return s
}

func (s Set[E]) Add(xs ...E) {
	for _, x := range xs {
		s[x] = keyExists
	}
}

func (s Set[E]) Remove(xs ...E) {
	for _, x := range xs {
		delete(s, x)
	}
}

func (s Set[E]) Len() int {
	return len(s)
}

func (s Set[E]) Contains(x E) bool {
	_, exists := s[x]
	return exists
}

func (s Set[E]) Equals(other Set[E]) bool {
	if s.Len() != other.Len() {
		return false
	}

	for x := range s {
		if !other.Contains(x) {
			return false
		}
	}

	return true
}

func (s Set[E]) Items() []E {
	res := make([]E, 0, s.Len())
	for x := range s {
		res = append(res, x)
	}
	return res
}

func (s Set[E]) Clone() Set[E] {
	res := NewSet[E]()
	for x := range s {
		res.Add(x)
	}
	return res
}

func (s Set[E]) Difference(other Set[E]) Set[E] {
	res := NewSet[E]()
	for x := range s {
		if !other.Contains(x) {
			res.Add(x)
		}
	}
	return res
}

func (s Set[E]) Intersection(other Set[E]) Set[E] {
	res := NewSet[E]()
	if s.Len() < other.Len() {
		for x := range s {
			if other.Contains(x) {
				res.Add(x)
			}
		}
	} else {
		for x := range other {
			if s.Contains(x) {
				res.Add(x)
			}
		}
	}
	return res
}

func (s Set[E]) Union(other Set[E]) Set[E] {
	res := s.Clone()
	for x := range other {
		res.Add(x)
	}
	return res
}
