# Go! Tiger Style

A set of ammendmentments to [Tigerstyle](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md),
for Golang.

First, I would read [TigerStyle](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md) fully.

This will provide a solid foundation for the Go! specific ammendements listed herein.

## Why Style?

As taken directly from the TigerStyle docs:

> "...in programming, style is not something to pursue directly. Style is necessary only where
> understanding is missing."

For our operations at Predixus, we are building data driven applications in Go. This, if done earnestly, results
in us going into the unknown and recovering the preverbial gold to distribute. And so, there

And so the goal of a style is to guide and support development with 3 north stars:

- Effectiveness - The software _has_ to work
- Simplicity - The design goals need to come together in as simple a manner as possible
- Robustness - There needs to be a baseline level of robustness

So, that being said, let's get into it.

## Technical Debt

There is nothing to add here. Tigerstyle nailed it. Refer to [their comments](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#technical-debt) on technical debt.

## Safety

[NASAs Power of 10](https://spinroot.com/gerard/pdf/P10.pdf) still applies. But there are some modifications
that we need to make that are specific to Go!:

- Use explicitely-sized types for everything:

  ```go
    type myVar int16
  ```

  instead of

  ```go
    type myVar int
  ```

- Specify capacity when allocating using `make`:

  ```go
    myVar := make([]int16, 5, 5)
    cap(myVar) // 5
  ```

  instead of

  ```go
    myVar := make([]int16, 5)
    cap(myVar) // 5
  ```

- Assertions:

  - Go does not have any natural notion of `assert`. The Go development team have stated their view on this:
    > ["...programmers use them as a crutch to avoid thinking about proper error handling and reporting"](https://go.dev/doc/faq#assertions)

  Completely transparent error handling in Go is indeed one of its strongest features - there is no need to use
  assertions to replace it. But there is value in using assertions to capture programmer errors that should be
  be caught during the testing phase. As stated in Tigerstyle:

  > "Assertions are a force multiplier for discovering bugs by fuzzing."

  To achieve this, build tags should be used to build release and debug assertion functions:

  ```go
  //assert_debug.go
  //go:build debug
  package main

  func assert(condition boolean, msg string) {
      if !condition {
          panic(msg)
      }
  }

  //assert_release.go
  //go:build !debug
  package main

  func assert(condition boolean, msg string) {
      if condition {
          panic(msg)
      }
  }
  ```

  And use like so:

  ```go
      package main

      type Counter struct {
          count     int16
      }

      func (c *Counter) Increment() {
          c.count += 1
      }

      func (c *Counter) Reset() {
          c.count = 0
      }

      func (c *Counter) Update(u int16) {
          c.count += u
          assert(c.count>=0, "Count cannot be negative")
      }

      func (c *Counter) Count() int16 {
          return c.count
      }
  ```

  If an assertion is raised, then the Counter has been incorrectly updated with a negative integer. This would highlight two programmer errors:

  1. `count` should be of type `int16`
  2. The calling code (incorrectly) expected to pass negative integers

- Assert the Positive Space, Negative Space and _Property Space_. In the plane of variable expressions, asserting
  just the positive and negative space Presents as points.
