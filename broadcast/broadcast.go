package broadcast

import "context"

// Broadcaster is a communication service with one sender and many recievers with all recievers (subscribers)
// getting every message sent by the sender. All communication happens via channels.
type Broadcaster[T any] struct {
	source         <-chan T
	listeners      map[chan T]struct{}
	addListener    chan chan T
	removeListener chan (<-chan T)
	bufferSize     int
	broadcast      func(context.Context, T)
}

// New creates a Broadcaster that will forward all data from the source channel to subscribers.
// All subscribers' channels will have the same capacity as source.
func New[T any](ctx context.Context, source <-chan T) *Broadcaster[T] {
	return NewBuffered(ctx, source, cap(source), cap(source) == 0)
}

// NewUnbuffered creates a Broadcaster that will forward all data from the source channel to subscribers.
// All subscribers' channels will be unbuffered (capacity = 0).
func NewUnbuffered[T any](ctx context.Context, source <-chan T) *Broadcaster[T] {
	return NewBuffered(ctx, source, 0, true)
}

// NewBuffered creates a Broadcaster that will forward all data from the source channel to subscribers.
// All subscribers' channels will have the capacity of bufferSize. Depending on the value of blockOnSend the Broadcaster
// will either ensure that all messages are sent to all subscribers (even if that means waiting on unbuffered/full channels)
// or it can skip any channel whose buffer is either full or 0.
func NewBuffered[T any](ctx context.Context, source <-chan T, bufferSize int, blockOnSend bool) *Broadcaster[T] {
	service := &Broadcaster[T]{
		source:         source,
		listeners:      make(map[chan T]struct{}),
		addListener:    make(chan chan T),
		removeListener: make(chan (<-chan T)),
		bufferSize:     bufferSize,
	}

	// Do we want broadcast to block on slow consumers? At least for bufferSize=0 we might have to.
	if blockOnSend {
		service.broadcast = service.ensureBroadcast // Blocking send. Ensures all subscribers get the data.
	} else {
		service.broadcast = service.tryBroadcast // Unblocking send. Doesn't send to a channel if that channel is full.
	}

	go service.serve(ctx)
	return service
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
// is the only one mutating internal state we don't need any locks or syncronisation
// other than using channels for communication.
func (s *Broadcaster[T]) serve(ctx context.Context) {
	defer s.closeAll()
	for {
		select {
		case <-ctx.Done():
			return
		case newListener := <-s.addListener:
			s.listeners[newListener] = struct{}{}
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
	// Cannot close using a read-only channel or remove the channel from a map.
	for listener := range s.listeners {
		if listener == sub {
			delete(s.listeners, listener)
			close(listener)
		}
	}
}

func (s *Broadcaster[T]) closeAll() {
	for listener := range s.listeners {
		close(listener)
	}
}

func (s *Broadcaster[T]) ensureBroadcast(ctx context.Context, val T) {
	for listener := range s.listeners {
		select {
		case listener <- val:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Broadcaster[T]) tryBroadcast(_ context.Context, val T) {
	for listener := range s.listeners {
		select {
		case listener <- val:
		default:
			// TODO: log? report? unsubscribe?
			s.removeSubscriber(listener)
		}
	}
}
