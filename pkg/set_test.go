package adt

import (
	"fmt"
	"reflect"
	"sort"
)

func ExampleNewSet() {
	s1 := NewSet(1, 2, 3)
	fmt.Println(reflect.TypeOf(s1))

	s2 := NewSet("a", "b", "c")
	fmt.Println(reflect.TypeOf(s2))

	s3 := NewSet[rune]() // example of initialization with 0 elements
	fmt.Println(reflect.TypeOf(s3))

	// Output:
	// adt.Set[int]
	// adt.Set[string]
	// adt.Set[int32]
}

func ExampleSet_Add() {
	s := NewSet(1)
	s.Add(2)
	s.Add(2, 3, 4)

	xs := s.Items()
	sort.Ints(xs)
	fmt.Println(xs)

	// Output: [1 2 3 4]
}

func ExampleSet_Remove() {
	s := NewSet(1, 2, 3, 4, 5)
	s.Remove(7, 6, 5, 1)

	xs := s.Items()
	sort.Ints(xs)
	fmt.Println(xs)

	// Output: [2 3 4]
}

func ExampleSet_Len() {
	s := NewSet[int]()
	fmt.Println(s.Len())

	s.Add(1, 2, 3)
	fmt.Println(s.Len())

	s.Remove(1, 2, 3)
	fmt.Println(s.Len())

	// Output:
	// 0
	// 3
	// 0
}

func ExampleSet_Contains() {
	s := NewSet[string]()
	fmt.Println(s.Contains("a"))

	s.Add("a")
	fmt.Println(s.Contains("a"))
	fmt.Println(s.Contains("b"))

	// Output:
	// false
	// true
	// false
}

func ExampleSet_Equals() {
	s1 := NewSet("a", "b", "c")
	s2 := NewSet("b", "c")
	fmt.Println(s1.Equals(s2))
	fmt.Println(s2.Equals(s1))

	s2.Add("a")
	fmt.Println(s1.Equals(s2))
	fmt.Println(s2.Equals(s1))

	// Output:
	// false
	// false
	// true
	// true
}

func ExampleSet_Items() {
	s := NewSet[int64]()
	fmt.Println(s.Items())

	s.Add(1)
	fmt.Println(s.Items())

	// Output:
	// []
	// [1]
}

func ExampleSet_Clone() {
	s1 := NewSet(1, 2, 3, 4)
	s2 := s1.Clone()

	fmt.Println(&s1 == &s2)
	fmt.Println(s1.Equals(s2))

	// Output:
	// false
	// true
}

func ExampleSet_Difference() {
	s1 := NewSet(1, 2, 3, 4, 5)
	s2 := NewSet(1, 2, 3)
	res1 := s1.Difference(s2).Items()
	res2 := s2.Difference(s1).Items()

	sort.Ints(res1)
	fmt.Println(res1)
	fmt.Println(res2)

	// Output:
	// [4 5]
	// []
}

func ExampleSet_Intersection() {
	s1 := NewSet(1, 2, 3, 4, 5)
	s2 := NewSet(1, 2, 3)
	res1 := s1.Intersection(s2).Items()
	res2 := s2.Intersection(s1).Items()

	sort.Ints(res1)
	sort.Ints(res2)
	fmt.Println(res1)
	fmt.Println(res2)

	// Output:
	// [1 2 3]
	// [1 2 3]
}

func ExampleSet_Union() {
	s1 := NewSet(1, 2, 3, 4, 5)
	s2 := NewSet(4, 5, 6, 7, 8)
	res1 := s1.Union(s2).Items()
	res2 := s2.Union(s1).Items()

	sort.Ints(res1)
	sort.Ints(res2)
	fmt.Println(res1)
	fmt.Println(res2)

	// Output:
	// [1 2 3 4 5 6 7 8]
	// [1 2 3 4 5 6 7 8]
}
