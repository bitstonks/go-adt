package broadcast

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleBroadcaster() {
	source := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := New(ctx, source)
	sub1 := broadcast.Subscribe()
	sub2 := broadcast.Subscribe()
	source <- 5318008
	fmt.Println("Sub1:", <-sub1)
	fmt.Println("Sub2:", <-sub2)
	// Output:
	// Sub1: 5318008
	// Sub2: 5318008
}

func TestBroadcaster_Subscribe(t *testing.T) {
	source := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := New(ctx, source)
	sub1 := broadcast.Subscribe()
	sub2 := broadcast.Subscribe()
	source <- 5318008
	assert.Equal(t, 5318008, <-sub1)
	assert.Equal(t, 5318008, <-sub2)

	source <- 42
	assert.Equal(t, 42, <-sub1)
	assert.Equal(t, 42, <-sub2)
}

func TestBroadcaster_Unsubscribe(t *testing.T) {
	source := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := New(ctx, source)
	sub1 := broadcast.Subscribe()
	sub2 := broadcast.Subscribe()

	broadcast.Unsubscribe(sub1)
	source <- 5318008
	_, ok := <-sub1
	assert.Equal(t, false, ok)
	assert.Equal(t, 5318008, <-sub2)

	broadcast.Unsubscribe(sub2)
	_, ok = <-sub2
	assert.Equal(t, false, ok)
}

func TestBufferedBroadcaster(t *testing.T) {
	source := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := NewBuffered(ctx, source, 10, false)
	sub := broadcast.Subscribe()

	// Overfill the buffer
	for i := 0; i < 20; i++ {
		source <- i
	}

	// Consume buffer that is only the first 10 elements
	for i := 0; i < 10; i++ {
		assert.Equal(t, i, <-sub)
	}
	assert.Len(t, sub, 0)
}

func TestBlockingBroadcaster(t *testing.T) {
	source := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := NewBuffered(ctx, source, 10, true)
	sub := broadcast.Subscribe()

	// (Over)Fill the buffer
	for i := 0; i < 11; i++ {
		source <- i
	}

	// Consume buffer
	for i := 0; i < 11; i++ {
		assert.Equal(t, i, <-sub)
	}

	// (Over)Fill the buffer
	for i := 0; i < 11; i++ {
		source <- i
	}

	canSend := true
	select {
	case source <- 11:
	default:
		canSend = false
	}
	assert.False(t, canSend)

	// Consume buffer
	for i := 0; i < 11; i++ {
		assert.Equal(t, i, <-sub)
	}
	assert.Len(t, sub, 0)
}

func TestUnbufferedBroadcaster(t *testing.T) {
	source := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := NewUnbuffered(ctx, source)
	sub := broadcast.Subscribe()

	source <- 5318008
	assert.Equal(t, 5318008, <-sub)

	source <- 5318008
	canSend := true
	select {
	case source <- 42:
	default:
		canSend = false
	}
	assert.False(t, canSend)
	assert.Equal(t, 5318008, <-sub)
}

func TestBroadcaster_closeAll(t *testing.T) {
	source := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	broadcast := New(ctx, source)
	sub1 := broadcast.Subscribe()
	sub2 := broadcast.Subscribe()
	source <- 5318008
	assert.Equal(t, 5318008, <-sub1)
	assert.Equal(t, 5318008, <-sub2)

	cancel()
	_, ok := <-sub1
	assert.Equal(t, false, ok)
	_, ok = <-sub2
	assert.Equal(t, false, ok)
}
