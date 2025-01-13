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

## Technical Debt

There is nothing to add here. Tigerstyle nailed it. Refer to
[their comments](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#technical-debt) on technical debt.

## Safety & Performance

The code in many of these ammenedements has been implemented, tested, benchmarked and fuzzed. Look at [bench.txt](bench.txt)
for details on the benchmarks and the `main.go` & `main_test.go` files for details on the tests.

[NASAs Power of 10](https://spinroot.com/gerard/pdf/P10.pdf) still applies. But there are some modifications
that we need to make that are specific to Go!:

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

2. Understanding Capacity Sharing Between Slices

   There are two approaches to handling slice capacity sharing, each with different trade-offs:

   - The first, is safe but costly (in allocations). This approach should be used when slice independence
     needs to be garunteed.

     ```go
     original := make([]int16, 5, 5)
     slice2 := make([]int16, 3)
     copy(slice2, original[0:3]) // Copy just the elements we want
     slice2 = append(slice2, 6)  // Now this truly won't affect original
     ```

   - The second approach is fast, but requires careful handling as capacity of a single slice is shared.
     Use when performance is critical and the implications are well understood.

     ```go
     original := make([]int16, 5, 10)
     slice2 := original[0:3]     // slice2 shares backing array
     slice2 = append(slice2, 6)  // Modifies original's backing array!
     ```

   Choose between these patterns based on your needs:

   - Use no-sharing when slice independence is crucial for correctness
   - Use sharing when performance is critical and you can carefully manage the slice relationships
   - The sharing approach is ~140x faster but requires more careful programming (inspect the benchmarks)

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

   Choose based on your requirements:

   - Use fixed-size pools when memory constraints are critical
   - Use growable pools when performance is the priority

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

- Assert the _Property Space_ wherever possible, and use Golangs Fuzzer to test it

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
  )

  // Traditional table test - tests input output pairs
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
  	seeds := []string{"", "a", "hello", "12345", "!@#$%"}
  	for _, seed := range seeds {
  		f.Add(seed)
  	}

  	f.Fuzz(func(t *testing.T, input string) {
  		// Property 1: reversing twice should return the original string
  		if twice := ReverseString(ReverseString(input)); twice != input {
  			t.Errorf("Double reverse failed: got %q, want %q", twice, input)
  		}

  		// Property 2: byte length should be preserved
  		reversed := ReverseString(input)
  		if len(reversed) != len(input) {
  			t.Errorf("Length not preserved: got %d bytes, want %d bytes",
  				len(reversed), len(input))
  		}
  	})
  }
  ```

  Key properties to consider testing when fuzzing:

  1. **Invariants**: Properties that should always hold true
  2. **Inverse operations**: Operations that should cancel each other out. E.g. encode/decode
  3. **Idempotency**: Operations that yield the same result when applied multiple times
  4. **Non-Idempotency**: Operations that do _not_ yield the same result when applied multiple
     times. E.g. a hashing algorithm

  When using runtime assertions on properties, focus on invariants that indicate
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
    strictest settings
