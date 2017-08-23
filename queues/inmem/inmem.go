package inmem

import (
	"sync"
)

var transport *Transport

// Transport is a vice.Transport for In-Memory queue
type Transport struct {
	m        sync.Mutex
	wg       sync.WaitGroup
	chans    map[string]chan []byte
	errChan  chan error
	stopchan chan struct{}
	stopped  bool
}

// New returns a new transport
func New() *Transport {
	if transport == nil {
		transport = &Transport{
			chans:    make(map[string]chan []byte),
			errChan:  make(chan error),
			stopchan: make(chan struct{}),
		}
	}
	return transport
}

// Receive gets a channel on which to receive messages
// with the specified name. The name is the queue's url
func (t *Transport) Receive(name string) <-chan []byte {
	t.m.Lock()
	defer t.m.Unlock()

	ch, ok := t.chans[name]
	if ok {
		return ch
	}
	return t.makeTopicChan(name)
}

func (t *Transport) makeTopicChan(name string) chan []byte {
	ch := make(chan []byte, 1024)
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for {
			select {
			case <-t.stopchan:
				if len(ch) != 0 {
					continue
				}
				close(ch)
				return
			default:
				continue
			}
		}
	}()
	t.chans[name] = ch
	return ch
}

// Send gets a channel on which messages with the
// specified name may be sent. The name is the queue's
// URL
func (t *Transport) Send(name string) chan<- []byte {
	t.m.Lock()
	defer t.m.Unlock()

	ch, ok := t.chans[name]
	if ok {
		return ch
	}
	return t.makeTopicChan(name)
}

// ErrChan gets the channel on which errors are sent.
func (t *Transport) ErrChan() <-chan error {
	return t.errChan
}

// Stop stops the transport.
// The channel returned from Done() will be closed
// when the transport has stopped.
func (t *Transport) Stop() {
	if t.stopped {
		return
	}
	t.wg.Wait()
	close(t.stopchan)
	close(t.errChan)
	t.stopped = true
}

// Done gets a channel which is closed when the
// transport has successfully stopped.
func (t *Transport) Done() chan struct{} {
	return t.stopchan
}
