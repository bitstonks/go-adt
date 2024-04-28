package deque

import (
	"testing"
)

var sink int

func BenchmarkDeque_PushBack(b *testing.B) {
	var d Deque[int]
	for i := range b.N {
		d.PushBack(i)
	}
}
func BenchmarkDeque_PushFront(b *testing.B) {
	var d Deque[int]
	for i := range b.N {
		d.PushFront(i)
	}
}
func BenchmarkDeque_PopBack(b *testing.B) {
	d := Deque[int]{data: make([]int, b.N), size: uint(b.N), first: 0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = d.PopBack()
	}
}
func BenchmarkDeque_PopFront(b *testing.B) {
	d := Deque[int]{data: make([]int, b.N), size: uint(b.N), first: 0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = d.PopFront()
	}
}
func BenchmarkDeque_Queue(b *testing.B) {
	q := New[int](2)
	q.PushBack(0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = q.PopFront()
		q.PushBack(0)
	}
}
func BenchmarkSlice_Queue(b *testing.B) {
	q := make([]int, 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = q[0]
		q = q[1:]
		q = append(q, 0)
	}
}
