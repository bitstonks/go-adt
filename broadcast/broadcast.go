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

// Broadcaster is a communication service with one sender and many recievers with all recievers (subscribers)
// getting every message sent by the sender. All communication happens via channels.
type Broadcaster[T any] struct {
	source         <-chan T
	listeners      map[<-chan T]chan T
	addListener    chan chan T
	removeListener chan (<-chan T)
	bufferSize     int
	broadcast      func(context.Context, T)
}

// NewBestEffort creates a Broadcaster that will try to forward data from the source channel to subscribers.
// It will skip any channels that are full i.e. the subscribers are too slow at emptying it.
// All subscribers' channels will have the same capacity as source or 1 in case of unbuffered source.
func NewBestEffort[T any](ctx context.Context, source <-chan T) *Broadcaster[T] {
	buf := cap(source)
	if buf == 0 {
		buf = 1
	}
	b, _ := New(ctx, source, buf, Skip) // We know this can't cause errors
	return b
}

// NewSynchronous creates a Broadcaster that will ensure all messages are not only sent, but delivered
// to subscribers.
func NewSynchronous[T any](ctx context.Context, source <-chan T) *Broadcaster[T] {
	b, _ := New(ctx, source, 0, Wait) // We know this can't cause errors
	return b
}

// New creates a Broadcaster that will forward all data from the source channel to subscribers.
// All subscribers' channels will have the capacity of bufferSize. Depending on the deliveryStrategy the Broadcaster
// will either
// * ensure that all messages are sent to all subscribers (even if that means waiting on unbuffered/full channels),
// * skip any channel whose buffer is full,
// * broadcast to empty channels and unsubscribe the rest.
func New[T any](ctx context.Context, source <-chan T, bufferSize int, deliveryStrategy DeliveryStrategy) (*Broadcaster[T], error) {
	if deliveryStrategy != Wait && bufferSize == 0 {
		return nil, fmt.Errorf("unbuffered channels only allowed for Wait strategy, not %q", deliveryStrategy)
	}
	service := &Broadcaster[T]{
		source:         source,
		listeners:      make(map[<-chan T]chan T),
		addListener:    make(chan chan T),
		removeListener: make(chan (<-chan T)),
		bufferSize:     bufferSize,
	}

	switch deliveryStrategy {
	case Skip:
		service.broadcast = service.broadcastOrSkip
	case Wait:
		service.broadcast = service.broadcastOrWait
	case Unsubscribe:
		service.broadcast = service.broadcastOrUnsubscribe
	default:
		return nil, fmt.Errorf("unknown value for deliveryStrategy: %q", deliveryStrategy)
	}

	go service.serve(ctx)
	return service, nil
}

// Subscribe will return a read-only channel that will deliver all broadcast messages to a new subscriber.
func (s *Broadcaster[T]) Subscribe() <-chan T {
	newListener := make(chan T, s.bufferSize)
	s.addListener <- newListener
	return newListener
}

// Unsubscribe will close the given channel and ensure it doesn't recieve any more updates.
func (s *Broadcaster[T]) Unsubscribe(channel <-chan T) {
	s.removeListener <- channel
}

// serve is the main event-handling loop. Because the goroutine running this method
// is the only one mutating internal state we don't need any locks or synchronization
// other than using channels for communication.
func (s *Broadcaster[T]) serve(ctx context.Context) {
	defer s.closeAll()
	for {
		select {
		case <-ctx.Done():
			return
		case newListener := <-s.addListener:
			s.listeners[newListener] = newListener
		case listenerToRemove := <-s.removeListener:
			s.removeSubscriber(listenerToRemove)
		case val, ok := <-s.source:
			if !ok { // Source channel was closed.
				return
			}
			s.broadcast(ctx, val)
		}
	}
}

func (s *Broadcaster[T]) removeSubscriber(sub <-chan T) {
	if listener, ok := s.listeners[sub]; ok {
		delete(s.listeners, sub)
		close(listener)
	}
}

func (s *Broadcaster[T]) closeAll() {
	for _, listener := range s.listeners {
		close(listener)
	}
}

func (s *Broadcaster[T]) broadcastOrWait(ctx context.Context, val T) {
	for _, listener := range s.listeners {
		select {
		case listener <- val:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Broadcaster[T]) broadcastOrUnsubscribe(_ context.Context, val T) {
	for _, listener := range s.listeners {
		select {
		case listener <- val:
		default:
			s.removeSubscriber(listener)
		}
	}
}

func (s *Broadcaster[T]) broadcastOrSkip(_ context.Context, val T) {
	for _, listener := range s.listeners {
		select {
		case listener <- val:
		default:
		}
	}
}
