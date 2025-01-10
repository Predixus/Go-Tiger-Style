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
in us going into the unknown and recovering the proverbial gold to distribute. And so the goal of a style is to
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

- Always be explicit about capacity, when allocating via `make`:

  1. Preventing Hidden Allocations in Slices

     Pre-allocate capacity when the final size is known to avoid expensive grow operations:

     ```go
     // Good: Single allocation with known final size
     data := make([]int16, 0, 1000)
     for i := 0; i < 1000; i++ {
         data = append(data, int16(i))  // No reallocation needed
     }

     // Bad: Multiple allocations as slice grows
     data := make([]int16, 0)
     for i := 0; i < 1000; i++ {
         data = append(data, int16(i))  // May trigger reallocation
     }
     ```

  2. Preventing Capacity Sharing Between Slices

     Explicitly match capacity to length when you want to prevent slice operations from
     accessing underlying array capacity:

     ```go
     // Good: No capacity sharing
     original := make([]int16, 5, 5)
     slice2 := original[0:3]     // slice2 has capacity of 3
     slice2 = append(slice2, 6)  // Forces new allocation, original unchanged

     // Bad: Hidden capacity sharing
     original := make([]int16, 5, 10)
     slice2 := original[0:3]     // slice2 has capacity of 7
     slice2 = append(slice2, 6)  // Modifies original's backing array!
     ```

  3. Preventing Rehashing when Initialising Maps

     Pre-size maps when the approximate size is known to avoid expensive rehashing operations:

     ```go
     // Good: Single hash table allocation
     users := make(map[string]User, 1000)
     for i := 0; i < 1000; i++ {
         users[fmt.Sprintf("user%d", i)] = User{}  // No rehashing needed
     }

     // Bad: Multiple rehashing operations
     users := make(map[string]User)  // Default small capacity
     for i := 0; i < 1000; i++ {
         users[fmt.Sprintf("user%d", i)] = User{}  // Forces periodic rehashing
     }
     ```

  4. Preventing Deadlocks in Channels

     Be explicit about channel buffering intent, in the name of the variable to
     prevent accidental deadlocks:

     ```go
     // Good: Clear buffering intent for synchronous communication
     chSync := make(chan int)

     // Good: Buffered for async communication, up to a capacity
     chAsync := make(chan int, 5)

     // Bad: Default to unbuffered without considering communication patterns
     ch := make(chan int)  // Might deadlock if async communication is needed
     ```

  5. Explicit size Buffer Pools to Prevent Growth

     When implementing buffer pools, explicit capacity helps prevent buffer growth:

     ```go
     // Good: Fixed-size buffer pool
     type Pool struct {
         buffers sync.Pool
     }

     func NewPool() *Pool {
         return &Pool{
             buffers: sync.Pool{
                 New: func() interface{} {
                     return make([]byte, 0, 4096)  // Fixed capacity
                 },
             },
         }
     }

     // Bad: Growable buffers can escape size constraints
     func NewPool() *Pool {
         return &Pool{
             buffers: sync.Pool{
                 New: func() interface{} {
                     return make([]byte, 0)  // Can grow unbounded
                 },
             },
         }
     }
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

  func assert(condition boolean, msg string) {}
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

  - Use Go's static analysis tools (`go vet`, `staticcheck`, `golangci-lint`) at their
    strictest settings.
