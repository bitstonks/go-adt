// deque provides a generic implementation of a double-ended queue.
//
// Using it as a queue is faster than using a slice, because it doesn't need to shift elements.
// Below is a benchmark result showing 1.7x speedup when using deque instead of a slice.
// ```
// $ GOMAXPROCS=1 go test -benchmem -run=^$ -bench ^Benchmark.+_Queue$ ./deque -count=3
// goos: darwin
// goarch: amd64
// pkg: github.com/bitstonks/go-adt/deque
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkDeque_Queue    71107292                16.84 ns/op            0 B/op          0 allocs/op
// BenchmarkDeque_Queue    70529484                17.62 ns/op            0 B/op          0 allocs/op
// BenchmarkDeque_Queue    70236950                16.94 ns/op            0 B/op          0 allocs/op
// BenchmarkSlice_Queue    40343143                29.30 ns/op           16 B/op          1 allocs/op
// BenchmarkSlice_Queue    41047531                29.22 ns/op           16 B/op          1 allocs/op
// BenchmarkSlice_Queue    40031304                28.95 ns/op           16 B/op          1 allocs/op
// PASS
// ok      github.com/bitstonks/go-adt/deque       7.508s
// ```
package deque

type Deque[T any] struct {
	data  []T
	first uint
	size  uint
}

// New creates a new deque.
func New[T any](capacity uint) Deque[T] {
	return Deque[T]{data: make([]T, capacity)}
}

// Len returns the number of elements in the deque.
func (d *Deque[T]) Len() uint {
	return d.size
}

// Cap returns the capacity of the deque.
func (d *Deque[T]) Cap() uint {
	return uint(len(d.data))
}

// PushBack adds an element to the back of the deque.
func (d *Deque[T]) PushBack(elems ...T) {
	n := uint(len(elems))
	if d.Len()+n > d.Cap() {
		d.ensureCapacity(n)
	}
	for i, elem := range elems {
		d.data[(d.first+d.size+uint(i))%d.Cap()] = elem
	}
	d.size += n
}

// PushFront adds an element to the front of the deque.
func (d *Deque[T]) PushFront(elems ...T) {
	n := uint(len(elems))
	if d.Len()+n > d.Cap() {
		d.ensureCapacity(n)
	}
	c := d.Cap()
	d.first = (d.first + c - n) % c
	for i, elem := range elems {
		d.data[(d.first+n-1-uint(i))%c] = elem
	}
	d.size += n
}

// PopBack removes and returns the element at the back of the deque.
func (d *Deque[T]) PopBack() T {
	if d.Len() == 0 {
		panic("deque: PopBack() called on empty deque")
	}
	d.size--
	idx := (d.first + d.size) % d.Cap()
	elem := d.data[idx]
	d.data[idx] = d.zeroElement()
	return elem
}

// PopFront removes and returns the element at the front of the deque.
func (d *Deque[T]) PopFront() T {
	if d.Len() == 0 {
		panic("deque: PopFront() called on empty deque")
	}
	elem := d.data[d.first]
	d.data[d.first] = d.zeroElement()
	d.first = (d.first + 1) % d.Cap()
	d.size--
	return elem
}

// Back returns the element at the back of the deque.
func (d *Deque[T]) Back() T {
	if d.Len() == 0 {
		panic("deque: Back() called on empty deque")
	}
	return d.data[(d.first+d.size-1)%d.Cap()]
}

// Front returns the element at the front of the deque.
func (d *Deque[T]) Front() T {
	if d.Len() == 0 {
		panic("deque: Front() called on empty deque")
	}
	return d.data[d.first]
}

// PopNBack removes and returns the last n elements from the deque.
func (d *Deque[T]) PopNBack(n uint) []T {
	if n > d.Len() {
		panic("deque: PopNBack() called with n > Len()")
	}
	elems := make([]T, 0, n)
	for i := d.first + d.size - 1; i >= d.first+d.size-n; i-- {
		elems = append(elems, d.data[i%d.Cap()])
		d.data[i%d.Cap()] = d.zeroElement()
	}
	d.size -= n
	return elems
}

// PopNFront removes and returns the last n elements from the deque.
func (d *Deque[T]) PopNFront(n uint) []T {
	if n > d.Len() {
		panic("deque: PopNFront() called with n > Len()")
	}
	c := d.Cap()
	elems := make([]T, 0, n)
	for i := d.first; i < d.first+n; i++ {
		elems = append(elems, d.data[i%c])
		d.data[i%c] = d.zeroElement()
	}
	d.first = (d.first + n) % c
	d.size -= n
	return elems
}

// checkCapacity will double the capacity of the deque until there is rome for n additional elements.
func (d *Deque[T]) ensureCapacity(n uint) {
	if d.Cap() == 0 {
		d.data = make([]T, n)
		return
	}
	for d.Len()+n > d.Cap() {
		d.data = append(d.data, d.data...) // Double the capacity.
	}
	// Zero all elements before the first element and after the last element.
	for i := d.first + d.size; i%d.Cap() != d.first; i++ {
		d.data[i%d.Cap()] = d.zeroElement()
	}
}

func (d *Deque[T]) zeroElement() (elem T) { return }
