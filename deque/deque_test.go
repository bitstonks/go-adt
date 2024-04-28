package deque

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleDeque_stack() {
	var stack Deque[int]
	stack.PushBack(1, 2)
	fmt.Println(stack.PopBack())
	stack.PushBack(3, 4, 5)
	fmt.Println(stack.PopBack())
	fmt.Println(stack.Back())
	fmt.Println(stack.PopNBack(2))
	fmt.Println(stack.PopBack())

	// Output:
	// 2
	// 5
	// 4
	// [4 3]
	// 1
}

func ExampleDeque_queue() {
	var queue Deque[int]
	queue.PushBack(1, 2)
	fmt.Println(queue.PopFront())
	queue.PushBack(3, 4, 5)
	fmt.Println(queue.PopFront())
	fmt.Println(queue.Front())
	fmt.Println(queue.PopNFront(2))
	fmt.Println(queue.PopFront())

	// Output:
	// 1
	// 2
	// 3
	// [3 4]
	// 5
}

func TestNew(t *testing.T) {
	t.Parallel()
	d := New[int](4)
	assert.NotNil(t, d)
	assert.Zero(t, d.Len())
	assert.EqualValues(t, 4, d.Cap())
}

func TestDeque_PushBack(t *testing.T) {
	t.Parallel()
	d := Deque[int]{data: make([]int, 2), first: 1, size: 0}
	d.PushBack(1)
	assert.EqualValues(t, []int{0, 1}, d.data)
	d.PushBack(2, 3)
	assert.EqualValues(t, []int{0, 1, 2, 3}, d.data)
}

func TestDeque_PushFront(t *testing.T) {
	t.Parallel()
	d := Deque[int]{data: make([]int, 2), first: 1, size: 0}
	d.PushFront(1)
	assert.EqualValues(t, []int{1, 0}, d.data)
	d.PushFront(2, 3)
	assert.EqualValues(t, []int{1, 0, 3, 2}, d.data)
}

func TestDeque_PopBack(t *testing.T) {
	t.Parallel()
	d := Deque[int]{data: []int{1, 2, 3, 4}, first: 2, size: 4}
	assert.Equal(t, 2, d.PopBack())
	assert.Equal(t, 1, d.PopBack())
	assert.Equal(t, 4, d.PopBack())
	assert.Equal(t, 3, d.PopBack())
	assert.PanicsWithValue(t, "deque: PopBack() called on empty deque", func() { d.PopBack() })
}
func TestDeque_PopFront(t *testing.T) {
	t.Parallel()
	d := Deque[int]{data: []int{1, 2, 3, 4}, first: 2, size: 4}
	assert.Equal(t, 3, d.PopFront())
	assert.Equal(t, 4, d.PopFront())
	assert.Equal(t, 1, d.PopFront())
	assert.Equal(t, 2, d.PopFront())
	assert.PanicsWithValue(t, "deque: PopFront() called on empty deque", func() { d.PopFront() })
}

func TestDeque_Back(t *testing.T) {
	t.Parallel()
	d := Deque[int]{data: []int{1, 2, 3, 4}, first: 2, size: 4}
	assert.Equal(t, 2, d.Back())
	d.PopBack()
	assert.Equal(t, 1, d.Back())
	d.PopNBack(3)
	assert.PanicsWithValue(t, "deque: Back() called on empty deque", func() { d.Back() })
}
func TestDeque_Front(t *testing.T) {
	t.Parallel()
	d := Deque[int]{data: []int{1, 2, 3, 4}, first: 2, size: 4}
	assert.Equal(t, 3, d.Front())
	d.PopFront()
	assert.Equal(t, 4, d.Front())
	d.PopNFront(3)
	assert.PanicsWithValue(t, "deque: Front() called on empty deque", func() { d.Front() })
}
func TestDeque_PopNBack(t *testing.T) {
	t.Parallel()
	d := Deque[int]{data: []int{1, 2, 3, 4}, first: 2, size: 4}
	assert.EqualValues(t, []int{2, 1}, d.PopNBack(2))
	assert.EqualValues(t, []int{0, 0, 3, 4}, d.data)
	assert.Equal(t, 4, d.Back())
	assert.PanicsWithValue(t, "deque: PopNBack() called with n > Len()", func() { d.PopNBack(3) })
}
func TestDeque_PopNFront(t *testing.T) {
	t.Parallel()
	d := Deque[int]{data: []int{1, 2, 3, 4}, first: 2, size: 4}
	assert.EqualValues(t, []int{3, 4}, d.PopNFront(2))
	assert.EqualValues(t, []int{1, 2, 0, 0}, d.data)
	assert.Equal(t, 1, d.Front())
	assert.PanicsWithValue(t, "deque: PopNFront() called with n > Len()", func() { d.PopNFront(3) })
}

func TestEnsureCapacity(t *testing.T) {
	t.Parallel()
	for _, capacity := range []uint{0, 1, 17} {
		capacity := capacity
		t.Run(fmt.Sprintf("from empty to %d", capacity), func(t *testing.T) {
			t.Parallel()
			var d Deque[int]
			d.ensureCapacity(capacity)
			assert.EqualValues(t, capacity, len(d.data))
			assert.Equal(t, capacity, d.Cap())
			assert.Zero(t, d.first)
			assert.Zero(t, d.size)
			for i := range capacity {
				assert.Zero(t, d.data[i])
			}
		})
	}
	t.Run("from one to 2", func(t *testing.T) {
		t.Parallel()
		d := Deque[int]{data: []int{1}, first: 0, size: 1}
		d.ensureCapacity(1)
		assert.EqualValues(t, []int{1, 0}, d.data)
	})
	t.Run("from one to 6", func(t *testing.T) {
		t.Parallel()
		d := Deque[int]{data: []int{1}, first: 0, size: 1}
		d.ensureCapacity(5)
		assert.EqualValues(t, []int{1, 0, 0, 0, 0, 0, 0, 0}, d.data)
	})
	t.Run("noop", func(t *testing.T) {
		t.Parallel()
		d := Deque[int]{data: []int{1, 2, 3, 0}, first: 0, size: 3}
		d.ensureCapacity(1)
		assert.EqualValues(t, []int{1, 2, 3, 0}, d.data)
	})
	t.Run("crossed", func(t *testing.T) {
		t.Parallel()
		d := Deque[int]{data: []int{3, 4, 1, 2}, first: 2, size: 4}
		d.ensureCapacity(1)
		assert.EqualValues(t, []int{0, 0, 1, 2, 3, 4, 0, 0}, d.data)
	})
}
