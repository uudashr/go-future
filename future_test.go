package future_test

import (
	"context"
	"testing"
	"time"

	future "github.com/uudashr/go-future"
)

func TestFuture_Get_immediate(t *testing.T) {
	fut, setResult := future.New()

	want := "Hello"
	setResult(want, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	got, err := fut.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatal("got:", got, "want:", want)
	}
}

func TestFuture_Listen_immediate(t *testing.T) {
	fut, setResult := future.New()

	want := "Hello"
	setResult(want, nil)

	done := make(chan struct{})
	var got string
	var gotErr error
	fut.Listen(func(val future.Value, err error) {
		got, gotErr = val.(string), err
		close(done)
	})

	<-done
	if gotErr != nil {
		t.Fatal(gotErr)
	}

	if got != want {
		t.Fatal("got:", got, "want:", want)
	}
}

func TestFuture_Get_async(t *testing.T) {
	fut, setResult := future.New()

	want := "Hello"
	time.AfterFunc(10*time.Millisecond, func() {
		setResult(want, nil)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	got, err := fut.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatal("got:", got, "want:", want)
	}
}

func TestFuture_Get_timeout(t *testing.T) {
	fut, _ := future.New()

	// no result

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := fut.Get(ctx)
	if got, want := err, context.DeadlineExceeded; got != want {
		t.Fatal("got:", got, "want:", want)
	}
}

func TestCall(t *testing.T) {
	want := "Hello"

	doThings := func() (string, error) {
		return want, nil
	}

	fut := future.Call(func() (future.Value, error) {
		return doThings()
	})

	got, err := fut.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatal("got:", got, "want:", want)
	}
}
