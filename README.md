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
in us going into the unknown and recovering the preverbial gold to distribute. And so the goal of a style is to
guide and support development through the unknown.

So, that being said, let's get into it.

## Technical Debt

There is nothing to add here. Tigerstyle nailed it. Refer to
[their comments](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#technical-debt) on technical debt.

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

  And used like so:

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

  If an assertion is raised, then the Counter has been incorrectly updated with a negative integer.
  This highlights two programmer errors:

  1. `count` should be of type `int16`
  2. The calling code expected to be able to pass negative integers

- Assert the _Property Space_ wherever possible.

  Property-based testing expands beyond traditional table-driven tests by verifying properties
  that should hold true for entire classes of inputs, rather than just specific examples. While
  table tests verify individual points in the input-output space, property tests verify relationships
  that should hold across the entire space.

  ![Property based test image](./assets/images/prop_based_test_dark.png#gh-dark-mode-only)
  ![Property based test image](./assets/images/prop_based_test_light.pngpng#gh-light-mode-only)

  For example, consider a function that reverses a string:

  ```go
  import (
      "testing"
      "github.com/stretchr/testify/assert"
  )

  // Traditional table test
  func TestReverse(t *testing.T) {
      tests := []struct {
          input    string
          expected string
      }{
          {"hello", "olleh"},
          {"world", "dlrow"},
      }
      for _, tt := range tests {
          got := Reverse(tt.input)
          assert.Equal(t, tt.expected, got, "Reverse(%q)", tt.input)
      }
  }

  // Property-based test
  func FuzzReverse(f *testing.F) {
    // Seed the corpus with the original test cases
    seeds := []string{"", "a", "hello", "12345", "!@#$%"}
    for _, seed := range seeds {
        f.Add(seed)
    }

    // Fuzz test that verifies two canoncial properties of `Reverse`
    f.Fuzz(func(t *testing.T, input string) {
        // Property 1: reversing twice should return the original string
        if twice := Reverse(Reverse(input)); twice!=input {
            t.Errorf("Double reverse failed: got %q, want %q", twice, input)
        }

        // Property 2: length should be preserved
        if reversed := Reverse(input); len(reversed) != len(input) {
            t.Errorf("Length not preserved: got %d, want %d", len(reversed), len(input))
        }
    })
  }
  ```

  Key properties to consider testing:

  1. **Invariants**: Properties that should always hold true
  2. **Inverse operations**: Operations that should cancel each other out. E.g. encode/decode
  3. **Idempotency**: Operations that yield the same result when applied multiple times
  4. **Non-Idempotency**: Operations that do _not_ yield the same result when applied multiple
     times. E.g. a hashing algorithm

  When using runtime assertions to test for properties, focus on invariants that indicate
  programmer errors:

  ```go
  func (v *Vector) Add(other *Vector) {
      assert(len(v.elements) == len(other.elements), "vectors must have same dimension")
      for i := range v.elements {
          v.elements[i] += other.elements[i]
      }
  }
  ```

  And remember: Property-based testing complements, not replaces, traditional testing
  approaches. Use both to achieve the required test coverage for your application.
