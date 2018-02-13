package future_test

import (
	"context"
	"fmt"
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

func ExampleFuture() {
	// Future function
	doThings := func() *future.Future {
		fut, setResult := future.New()
		time.AfterFunc(10*time.Millisecond, func() {
			setResult("OK", nil)
		})
		return fut
	}

	// Usage
	res := doThings()
	val, _ := res.Get(context.Background())
	fmt.Println(val)
	// Output: OK
}

func ExampleFuture_callback() {
	// Future function
	doThings := func() *future.Future {
		fut, setResult := future.New()
		time.AfterFunc(10*time.Millisecond, func() {
			setResult("OK", nil)
		})
		return fut
	}

	// Usage
	done := make(chan struct{})

	res := doThings()
	res.Listen(func(val future.Value, err error) {
		fmt.Println(val)
		close(done)
	})

	<-done
	// Output: OK
}

func ExampleCall() {
	// Sync function
	greet := func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "Hello World!", nil
	}

	// Use sync function as async using future
	fut := future.Call(func() (future.Value, error) {
		return greet()
	})

	v, _ := fut.Get(context.Background())
	fmt.Println(v)
	// Output: Hello World!

}
