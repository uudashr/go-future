[![GoDoc](https://godoc.org/github.com/uudashr/go-future?status.svg)](https://godoc.org/github.com/uudashr/go-future)
# Future

Future API in Golang. It give flexibility to the user to use both sync and async way to get the result.

The API is too generic. It use `type Value interface{}`, as we know `interface{}` says nothing. You can use this API or create specific implementation for each data type.

## Usage

Async function

```go
package main

func main() {
    res := doThings()
    val, _ := res.Get(context.Background())
    fmt.Println(val)
}

func doThings() *future.Future {
    fut, setResult := future.New()
    time.AfterFunc(10*time.Millisecond, func() {
        setResult("OK", nil)
    })
    return fut
}
```



Some of the function created in async manner, we can convert it as async using the future API.

```go
package main

func main() {
    // Use sync function as async using future
    fut := future.Call(func() (future.Value, error) {
		return greet()
    })

    v, _ := fut.Get(context.Background())
    fmt.Println(v)
}

// Sync function
func greet() (string, error) {
    time.Sleep(500 * time.Millisecond)
    return "Hello World!", nil
}
```



Use callback style to get the result immediately without the overhead of an extra goroutine

```go
package main

func main() {
    done := make(chan struct{})

    res := doThings()
    res.Listen(func(val future.Value, err error) {
        fmt.Println(val)
        close(done)
    })

    <-done
}

// Future function
func doThings() *future.Future {
    fut, setResult := future.New()
    time.AfterFunc(10*time.Millisecond, func() {
        setResult("OK", nil)
    })
    return fut
}
```

