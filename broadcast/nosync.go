package broadcast

import (
	"context"
)

// NoSyncBroadcaster is a broadcast service that allows a single sender to send
// messages to multiple recievers. It does not implement any synchronisation
// and can only be used if the sender is externally synchronised
// e.g. only used in one goroutine or using a mutex.
type NoSyncBroadcaster[T any] struct {
	subscribers map[<-chan T]chan T
	bufferSize  int
}

// NewNoSyncBroadcaster creates a NoSyncBroadcaster where all subscribers will
// get a channel of capacity bufferSize when they subscribe.
func NewNoSyncBroadcaster[T any](bufferSize int) *NoSyncBroadcaster[T] {
	return &NoSyncBroadcaster[T]{
		subscribers: make(map[<-chan T]chan T),
		bufferSize:  bufferSize,
	}
}

// Subscribe creates and returns a new channel that will receive all messages
// send by the sender via this broadcast service.
func (b *NoSyncBroadcaster[T]) Subscribe() <-chan T {
	ch := make(chan T, b.bufferSize)
	b.AddSubscriber(ch)
	return ch
}

// AddSubscriber gives subscribers the option to provide their own chanel to
// receive updates on. In case they already have allocated one and want to
// reuse it or if the default bufferSize isn't OK for them.
func (b *NoSyncBroadcaster[T]) AddSubscriber(sub chan T) {
	b.subscribers[sub] = sub
}

// Unsubscribe will stop the service sending messages on this channel and close
// the channel. Returns true if the provided channel is a valid subscriber.
func (b *NoSyncBroadcaster[T]) Unsubscribe(sub <-chan T) bool {
	if sub, ok := b.subscribers[sub]; ok {
		close(sub)
		delete(b.subscribers, sub)
		return true
	}
	return false
}

// CloseAll will close and delete all subscribers' channels.
func (b *NoSyncBroadcaster[T]) CloseAll() {
	for _, sub := range b.subscribers {
		close(sub)
		delete(b.subscribers, sub)
	}
}

// Len returns the number of subcribers that the service is send to.
func (b *NoSyncBroadcaster[T]) Len() int {
	return len(b.subscribers)
}

// SendOrWait will send message to all subscribers. If a subscriber's channel
// is full or unbuffered it will wait until a space in the channel frees up.
// It returns true if it managed to send the message to all subscribers before
// ctx expires.
func (b *NoSyncBroadcaster[T]) SendOrWait(ctx context.Context, message T) bool {
	for _, sub := range b.subscribers {
		select {
		case <-ctx.Done():
			return false
		case sub <- message:
		}
	}
	return true
}

// SendOrSkip will try to send message to all subscribers. If a subscriber's
// channel is full or unbuffered it will skip that channel and continue to the
// next one.
func (b *NoSyncBroadcaster[T]) SendOrSkip(message T) int {
	numSkipped := 0
	for _, sub := range b.subscribers {
		select {
		case sub <- message:
		default:
			numSkipped++
		}
	}
	return numSkipped
}

// SendOrUnsubscribe will try to send message to all subscribers. If a
// subscriber's channel is full or unbuffered it will unsubscribe that
// channel from further updates and close it.
func (b *NoSyncBroadcaster[T]) SendOrUnsubscribe(message T) int {
	numUnsub := 0
	for _, sub := range b.subscribers {
		select {
		case sub <- message:
		default:
			close(sub)
			delete(b.subscribers, sub)
			numUnsub++
		}
	}
	return numUnsub
}
