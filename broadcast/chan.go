package broadcast

import (
	"context"
	"fmt"
)

type DeliveryStrategy int

const (
	// Will try to push to channel and pass if channel is full.
	Skip DeliveryStrategy = iota
	// Will try to push to channel and wait if channel is full.
	Wait
	// Will try to push to channel and unsubscribe if channel is full.
	Unsubscribe
)

// ChanBroadcaster is a communication service with one sender and many recievers with all recievers (subscribers)
// getting every message sent by the sender. All communication happens via channels.
type ChanBroadcaster[T any] struct {
	source         <-chan T
	sender         *NoSyncBroadcaster[T]
	addListener    chan chan T
	removeListener chan (<-chan T)
	broadcast      func(context.Context, T)
}

// NewBestEffortChanBroadcaster creates a Broadcaster that will try to forward data from the source channel to subscribers.
// It will skip any channels that are full i.e. the subscribers are too slow at emptying it.
// All subscribers' channels will have the same capacity as source or 1 in case of unbuffered source.
func NewBestEffortChanBroadcaster[T any](ctx context.Context, source <-chan T) *ChanBroadcaster[T] {
	buf := cap(source)
	if buf == 0 {
		buf = 1
	}
	b, _ := NewChanBroadcaster(ctx, source, buf, Skip) // We know this can't cause errors
	return b
}

// NewSynchronousChanBroadcaster creates a Broadcaster that will ensure all messages are not only sent, but delivered
// to subscribers.
func NewSynchronousChanBroadcaster[T any](ctx context.Context, source <-chan T) *ChanBroadcaster[T] {
	b, _ := NewChanBroadcaster(ctx, source, 0, Wait) // We know this can't cause errors
	return b
}

// NewChanBroadcaster creates a Broadcaster that will forward all data from the source channel to subscribers.
// All subscribers' channels will have the capacity of bufferSize. Depending on the deliveryStrategy the Broadcaster
// will either
// * ensure that all messages are sent to all subscribers (even if that means waiting on unbuffered/full channels),
// * skip any channel whose buffer is full,
// * broadcast to empty channels and unsubscribe the rest.
func NewChanBroadcaster[T any](ctx context.Context, source <-chan T, bufferSize int, deliveryStrategy DeliveryStrategy) (*ChanBroadcaster[T], error) {
	if deliveryStrategy != Wait && bufferSize == 0 {
		return nil, fmt.Errorf("unbuffered channels only allowed for Wait strategy, not %q", deliveryStrategy)
	}
	service := &ChanBroadcaster[T]{
		source:         source,
		sender:         NewNoSyncBroadcaster[T](bufferSize),
		addListener:    make(chan chan T),
		removeListener: make(chan (<-chan T)),
	}

	switch deliveryStrategy {
	case Skip:
		service.broadcast = func(_ context.Context, message T) { service.sender.SendOrSkip(message) }
	case Wait:
		service.broadcast = func(ctx context.Context, message T) { service.sender.SendOrWait(ctx, message) }
	case Unsubscribe:
		service.broadcast = func(_ context.Context, message T) { service.sender.SendOrUnsubscribe(message) }
	default:
		return nil, fmt.Errorf("unknown value for deliveryStrategy: %q", deliveryStrategy)
	}

	go service.serve(ctx)
	return service, nil
}

// Subscribe will return a read-only channel that will deliver all broadcast messages to a new subscriber.
func (s *ChanBroadcaster[T]) Subscribe() <-chan T {
	newListener := make(chan T, s.sender.bufferSize)
	s.addListener <- newListener
	return newListener
}

// Unsubscribe will close the given channel and ensure it doesn't recieve any more updates.
func (s *ChanBroadcaster[T]) Unsubscribe(channel <-chan T) {
	s.removeListener <- channel
}

// serve is the main event-handling loop. Because the goroutine running this method
// is the only one mutating internal state we don't need any locks or synchronization
// other than using channels for communication.
func (s *ChanBroadcaster[T]) serve(ctx context.Context) {
	defer s.sender.CloseAll()
	for {
		select {
		case <-ctx.Done():
			return
		case newListener := <-s.addListener:
			s.sender.AddSubscriber(newListener)
		case listenerToRemove := <-s.removeListener:
			s.sender.Unsubscribe(listenerToRemove)
		case val, ok := <-s.source:
			if !ok { // Source channel was closed.
				return
			}
			s.broadcast(ctx, val)
		}
	}
}
