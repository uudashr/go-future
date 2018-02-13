package future

import (
	"context"
	"sync"
)

// Value represents value concept. Can be anything.
type Value interface{}

// SetResultFunc is function to set result of the Future.
type SetResultFunc func(Value, error)

// Future holds the value of Future.
type Future struct {
	val   Value
	err   error
	ready chan struct{}

	mu        sync.Mutex
	callbacks []SetResultFunc
}

// New constructs new Future.
func New() (*Future, SetResultFunc) {
	f := &Future{ready: make(chan struct{})}
	return f, f.setResult
}

// Get returns value when it's ready. Will return error when the ctx signal a cancelation.
func (f *Future) Get(ctx context.Context) (Value, error) {
	select {
	case <-f.ready:
		return f.val, f.err
	default:
	}

	select {
	case <-f.ready:
		return f.val, f.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// Ready indicates whether result ready or not.
func (f *Future) Ready() <-chan struct{} {
	return f.ready
}

func (f *Future) setResult(v Value, err error) {
	select {
	case <-f.ready:
	default:
		f.val, f.err = v, err
		f.notifyCallbacks()
	}
}

// Listen for the result.
func (f *Future) Listen(callback SetResultFunc) {
	select {
	case <-f.ready:
		callback(f.val, f.err)
	default:
		f.mu.Lock()
		f.callbacks = append(f.callbacks, callback)
		f.mu.Unlock()
	}
}

func (f *Future) notifyCallbacks() {

	close(f.ready)
	f.mu.Lock()
	for _, callback := range f.callbacks {
		callback(f.val, f.err)
	}
	f.mu.Unlock()
}

// Call will converts the sync function call as async call.
func Call(f func() (Value, error)) *Future {
	fut, setDone := New()
	go func() {
		res, err := f()
		setDone(res, err)
	}()
	return fut
}
