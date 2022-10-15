package broadcast

import (
	"context"
	"sync"
)

// SyncBroadcaster is a wrapper around NoSyncBroadcaster ensuring that all operations are
// properly synchronised using an internal mutex.
type SyncBroadcaster[T any] struct {
	nosync *NoSyncBroadcaster[T]
	lock   sync.RWMutex
}

// NewSyncBroadcaster creates a SyncBroadcaster where all subscribers will
// get a channel of capacity bufferSize when they subscribe.
func NewSyncBroadcaster[T any](bufferSize int) *SyncBroadcaster[T] {
	return &SyncBroadcaster[T]{
		nosync: NewNoSyncBroadcaster[T](bufferSize),
	}
}

// Subscribe creates and returns a new channel that will receive all messages
// send by the sender via this broadcast service.
func (b *SyncBroadcaster[T]) Subscribe() <-chan T {
	ch := make(chan T, b.nosync.bufferSize)
	b.AddSubscriber(ch)
	return ch
}

// AddSubscriber gives subscribers the option to provide their own chanel to
// receive updates on. In case they already have allocated one and want to
// reuse it or if the default bufferSize isn't OK for them.
func (b *SyncBroadcaster[T]) AddSubscriber(sub chan T) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.nosync.AddSubscriber(sub)
}

// Unsubscribe will stop the service sending messages on this channel and close
// the channel. Returns true if the provided channel is a valid subscriber.
func (b *SyncBroadcaster[T]) Unsubscribe(sub <-chan T) bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.nosync.Unsubscribe(sub)
}

// CloseAll will close and delete all subscribers' channels.
func (b *SyncBroadcaster[T]) CloseAll() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.nosync.CloseAll()
}

// Len returns the number of subcribers that the service is send to.
func (b *SyncBroadcaster[T]) Len() int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.nosync.Len()
}

// SendOrWait will send message to all subscribers. If a subscriber's channel
// is full or unbuffered it will wait until a space in the channel frees up.
// It returns true if it managed to send the message to all subscribers before
// ctx expires.
func (b *SyncBroadcaster[T]) SendOrWait(ctx context.Context, message T) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.nosync.SendOrWait(ctx, message)
}

// SendOrSkip will try to send message to all subscribers. If a subscriber's
// channel is full or unbuffered it will skip that channel and continue to the
// next one.
func (b *SyncBroadcaster[T]) SendOrSkip(message T) int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.nosync.SendOrSkip(message)
}

// SendOrUnsubscribe will try to send message to all subscribers. If a
// subscriber's channel is full or unbuffered it will unsubscribe that
// channel from further updates and close it.
func (b *SyncBroadcaster[T]) SendOrUnsubscribe(message T) int {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.nosync.SendOrUnsubscribe(message)
}
