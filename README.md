# clock [![GoDoc](https://godoc.org/github.com/tilinna/clock?status.png)](https://godoc.org/github.com/tilinna/clock)

A Go (golang) library for mocking standard time, optionally also with context.Context.

## USAGE

```go
m := clock.NewMock(time.Date(2018, 1, 1, 10, 0, 0, 0, time.UTC))
fmt.Println("Time is now", m.Now())
timer := m.NewTimer(15 * time.Second)
m.Add(25 * time.Second)
fmt.Println("Time is now", m.Now())
fmt.Println("Timeout was", <-timer.C)
// Output:
// Time is now 2018-01-01 10:00:00 +0000 UTC
// Time is now 2018-01-01 10:00:25 +0000 UTC
// Timeout was 2018-01-01 10:00:15 +0000 UTC
```

## TODO

- More tests
- Documentation with examples
- Tag v1.0
