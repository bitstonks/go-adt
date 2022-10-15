package broadcast

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleChanBroadcaster() {
	source := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := NewBestEffortChanBroadcaster(ctx, source)
	sub1 := broadcast.Subscribe()
	sub2 := broadcast.Subscribe()
	source <- 5318008
	fmt.Println("Sub1:", <-sub1)
	fmt.Println("Sub2:", <-sub2)
	// Output:
	// Sub1: 5318008
	// Sub2: 5318008
}

func TestChanBroadcaster_Subscribe(t *testing.T) {
	source := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := NewBestEffortChanBroadcaster(ctx, source)
	sub1 := broadcast.Subscribe()
	sub2 := broadcast.Subscribe()
	source <- 5318008
	assert.Equal(t, 5318008, <-sub1)
	assert.Equal(t, 5318008, <-sub2)

	source <- 42
	assert.Equal(t, 42, <-sub1)
	assert.Equal(t, 42, <-sub2)
}

func TestChanBroadcaster_Unsubscribe(t *testing.T) {
	source := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := NewBestEffortChanBroadcaster(ctx, source)
	sub1 := broadcast.Subscribe()
	sub2 := broadcast.Subscribe()

	// Ensure double Unsubscribe doesn't panic.
	broadcast.Unsubscribe(sub1)
	broadcast.Unsubscribe(sub1)
	source <- 5318008
	_, ok := <-sub1
	assert.Equal(t, false, ok)
	assert.Equal(t, 5318008, <-sub2)

	broadcast.Unsubscribe(sub2)
	_, ok = <-sub2
	assert.Equal(t, false, ok)
}

func TestUnsubscribingChanBroadcaster(t *testing.T) {
	source := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast, err := NewChanBroadcaster(ctx, source, 10, Unsubscribe)
	assert.NoError(t, err)
	sub := broadcast.Subscribe()

	// Overfill the buffer
	for i := 0; i < 20; i++ {
		source <- i
	}

	// Consume buffer that is only the first 10 elements
	for i := 0; i < 10; i++ {
		assert.Equal(t, i, <-sub)
	}
	// Channel is closed, becuase we couldn't keep up
	_, ok := <-sub
	assert.False(t, ok)
}

func TestWaitingChanBroadcaster(t *testing.T) {
	source := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast, err := NewChanBroadcaster(ctx, source, 10, Wait)
	assert.NoError(t, err)
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

func TestSynchronousChanBroadcaster(t *testing.T) {
	source := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broadcast := NewSynchronousChanBroadcaster(ctx, source)
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

func TestChanBroadcaster_closeAll(t *testing.T) {
	source := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	broadcast := NewBestEffortChanBroadcaster(ctx, source)
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

func testMultipleListeners(broadcast *ChanBroadcaster[int], source chan int, m int) func(*testing.B) {
	return func(b *testing.B) {
		var wg sync.WaitGroup
		wg.Add(m)
		for i := 0; i < m; i++ {
			sub := broadcast.Subscribe()
			go consume(sub, b.N, &wg)
			defer broadcast.Unsubscribe(sub)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			source <- i
		}
		wg.Wait()
	}
}

func BenchmarkWaitingChanBroadcaster(b *testing.B) {
	for _, bufSize := range []int{0, 1, 5, 10, 100} {
		source := make(chan int, bufSize)
		ctx, cancel := context.WithCancel(context.Background())
		broadcast, _ := NewChanBroadcaster(ctx, source, bufSize, Wait)
		for _, m := range []int{1, 2, 5, 10, 100} {
			b.Run(fmt.Sprintf("buffer:%v, subs:%v", bufSize, m), testMultipleListeners(broadcast, source, m))
		}
		cancel()
		close(source)
	}
}
