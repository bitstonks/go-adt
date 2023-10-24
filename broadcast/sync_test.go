package broadcast

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleSyncBroadcaster() {
	broadcast := NewSyncBroadcaster[int](10)
	sub := broadcast.Subscribe()
	broadcast.SendOrSkip(1)
	broadcast.SendOrSkip(3)
	fmt.Println(<-sub)
	fmt.Println(<-sub)
	broadcast.Unsubscribe(sub)
	if _, ok := <-sub; !ok {
		fmt.Println("Channel was closed and unsubscribed.")
	}
	// Output:
	// 1
	// 3
	// Channel was closed and unsubscribed.
}

func TestSyncBroadcaster_CloseAll(t *testing.T) {
	broadcast := NewSyncBroadcaster[int](10)
	var wg sync.WaitGroup
	n, m := 100, 10
	wg.Add(m)
	for i := 0; i < m; i++ {
		sub := broadcast.Subscribe()
		go consume(sub, n, &wg)
		defer func() {
			_, ok := <-sub
			assert.False(t, ok)
		}()
	}
	ctx := context.Background()
	for i := 0; i < n; i++ {
		broadcast.SendOrWait(ctx, i)
	}
	wg.Wait()
	broadcast.CloseAll()
	assert.Zero(t, broadcast.Len())
}

func TestSyncBroadcaster_SubUnsub(t *testing.T) {
	broadcast := NewSyncBroadcaster[int](10)
	var wg sync.WaitGroup
	n, m := 100, 10
	wg.Add(m)
	for i := 0; i < m; i++ {
		sub := broadcast.Subscribe()
		go consume(sub, n, &wg)
		defer func() {
			_, ok := <-sub
			assert.False(t, ok)
		}()
		defer broadcast.Unsubscribe(sub)
	}
	ctx := context.Background()
	for i := 0; i < n; i++ {
		broadcast.SendOrWait(ctx, i)
	}
	wg.Wait()
}

func TestSyncBroadcaster_SendOrWait(t *testing.T) {
	broadcast := NewSyncBroadcaster[int](1)
	sub := broadcast.Subscribe()
	ctx, cancel := context.WithCancel(context.Background())

	assert.Equal(t, true, broadcast.SendOrWait(ctx, 1))
	go func() { assert.Equal(t, true, broadcast.SendOrWait(ctx, 2)) }()
	assert.Equal(t, 1, <-sub)
	assert.Equal(t, 2, <-sub)

	assert.Equal(t, true, broadcast.SendOrWait(ctx, 3))
	cancel()
	assert.Equal(t, false, broadcast.SendOrWait(ctx, 4))
	broadcast.Unsubscribe(sub)
	assert.Equal(t, 3, <-sub)
	_, ok := <-sub
	assert.False(t, ok)
	assert.Zero(t, broadcast.Len())
}

func TestSyncBroadcaster_SendOrSkip(t *testing.T) {
	broadcast := NewSyncBroadcaster[int](1)
	sub := broadcast.Subscribe()
	assert.Equal(t, 0, broadcast.SendOrSkip(1))
	assert.Equal(t, 1, broadcast.SendOrSkip(2))
	assert.Equal(t, 1, <-sub)
	assert.Equal(t, 0, broadcast.SendOrSkip(3))
	assert.Equal(t, 3, <-sub)
	broadcast.Unsubscribe(sub)
	_, ok := <-sub
	assert.False(t, ok)
	assert.Zero(t, broadcast.Len())
}

func TestSyncBroadcaster_SendOrUnsubscribe(t *testing.T) {
	broadcast := NewSyncBroadcaster[int](1)
	sub := broadcast.Subscribe()
	assert.Equal(t, 0, broadcast.SendOrUnsubscribe(1))
	assert.Equal(t, 1, broadcast.SendOrUnsubscribe(2))
	assert.Equal(t, 1, <-sub)
	_, ok := <-sub
	assert.False(t, ok)
	assert.Zero(t, broadcast.Len())
}
