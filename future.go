package future

import "context"

// Value represents value concept. Can be anything.
type Value interface{}

// SetResultFunc is function to set result of the Future.
type SetResultFunc func(Value, error)

// Future holds the value of Future.
type Future struct {
	val       Value
	err       error
	ready     chan struct{}
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
		close(f.ready)
		f.notifyCallbacks()
	}
}

// Listen for the result.
func (f *Future) Listen(callback SetResultFunc) {
	select {
	case <-f.ready:
		callback(f.val, f.err)
	default:
		f.callbacks = append(f.callbacks, callback)
	}
}

func (f *Future) notifyCallbacks() {
	for _, callback := range f.callbacks {
		callback(f.val, f.err)
	}
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
