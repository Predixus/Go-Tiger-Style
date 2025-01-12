package main

import (
	"sync"
	"testing"
)

func BenchmarkSlice(b *testing.B) {
	b.Run("SliceAllocateCapacity", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			SliceAllocateCapacity()
		}
	})
	b.Run("SliceLetCapacityGrow", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			SliceLetCapacityGrow()
		}
	})
	b.Run("NoSliceCapacitySharing", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			SliceNoShareCapacity()
		}
	})

	b.Run("SliceCapacitySharing", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			SliceShareCapacity()
		}
	})
}

func TestSliceCapacityBehavior(t *testing.T) {
	t.Run("NoShareCapacity", func(t *testing.T) {
		original := SliceNoShareCapacity()

		if len(original) != 5 {
			t.Errorf("Original slice length changed: got %d, want 5", len(original))
		}
		if original[3] != 0 { // Now this should pass
			t.Errorf("Original slice was modified: got %d at index 3, want 0", original[3])
		}
	})

	t.Run("ShareCapacity", func(t *testing.T) {
		original := SliceShareCapacity()

		if original[3] != 6 {
			t.Errorf("Original slice was not modified: got %d at index 3, want 6", original[3])
		}
	})
}

// Map operations
func BenchmarkMap(b *testing.B) {
	b.Run("SingleHash", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			MapSingleHashAllocation()
		}
	})
	b.Run("MultipleHash", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			MapMultipleRehashings()
		}
	})
}

// Channel operations
func BenchmarkChannel(b *testing.B) {
	b.Run("Sync", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			ChannelSync()
		}
	})
	b.Run("Async", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			ChannelAsync()
		}
	})
	b.Run("SyncMulti", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			ChannelSyncMulti()
		}
	})
	b.Run("AsyncMulti", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			ChannelAsyncMulti()
		}
	})
}

// Test functions to verify behavior
func TestChannel(t *testing.T) {
	t.Run("Sync", func(t *testing.T) {
		val := ChannelSync()
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}
	})
	t.Run("Async", func(t *testing.T) {
		val := ChannelAsync()
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}
	})
	t.Run("SyncMulti", func(t *testing.T) {
		result := ChannelSyncMulti()
		if len(result) != 100 {
			t.Errorf("Expected length 100, got %d", len(result))
		}
		// Verify sequence
		for i := 0; i < 100; i++ {
			if result[i] != i {
				t.Errorf("At index %d: expected %d, got %d", i, i, result[i])
			}
		}
	})
	t.Run("AsyncMulti", func(t *testing.T) {
		result := ChannelAsyncMulti()
		if len(result) != 100 {
			t.Errorf("Expected length 100, got %d", len(result))
		}
		// Verify sequence
		for i := 0; i < 100; i++ {
			if result[i] != i {
				t.Errorf("At index %d: expected %d, got %d", i, i, result[i])
			}
		}
	})
}

func TestBufferPools(t *testing.T) {
	t.Run("FixedPool", func(t *testing.T) {
		pool := NewFixedPool(64)

		// Get a buffer and verify its capacity
		buf := pool.Get()
		if cap(buf) != 64 {
			t.Errorf("Expected capacity 64, got %d", cap(buf))
		}

		// Try to grow the buffer
		largerData := make([]byte, 128)
		grown := append(buf, largerData...)

		// Put back the grown buffer
		pool.Put(grown)

		// Get another buffer - should get a fresh one of original size
		newBuf := pool.Get()
		if cap(newBuf) != 64 {
			t.Errorf("Expected capacity 64 after grow attempt, got %d", cap(newBuf))
		}
	})

	t.Run("GrowablePool", func(t *testing.T) {
		pool := NewGrowablePool()

		// Get a buffer and grow it
		buf := pool.Get()
		largerData := make([]byte, 128)
		grown := append(buf, largerData...)

		// Put back the grown buffer
		pool.Put(grown)

		// Get another buffer - might get the grown one
		newBuf := pool.Get()
		if cap(newBuf) < 128 {
			t.Logf("Note: Got a fresh buffer instead of reusing grown one")
		}
	})
}

// Benchmark functions
func BenchmarkFixedPool(b *testing.B) {
	pool := NewFixedPool(64)
	data := make([]byte, 32) // Data smaller than buffer

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("Normal Use", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			buf := pool.Get()
			buf = append(buf, data...)
			pool.Put(buf)
		}
	})

	b.Run("Growth Attempt", func(b *testing.B) {
		largeData := make([]byte, 128) // Data larger than buffer
		for n := 0; n < b.N; n++ {
			buf := pool.Get()
			buf = append(buf, largeData...)
			pool.Put(buf) // Will discard grown buffer
		}
	})
}

func BenchmarkGrowablePool(b *testing.B) {
	pool := NewGrowablePool()
	data := make([]byte, 32)

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("Normal Use", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			buf := pool.Get()
			buf = append(buf, data...)
			pool.Put(buf)
		}
	})

	b.Run("Growth Allowed", func(b *testing.B) {
		largeData := make([]byte, 128)
		for n := 0; n < b.N; n++ {
			buf := pool.Get()
			buf = append(buf, largeData...)
			pool.Put(buf) // Will retain grown buffer
		}
	})
}

// Concurrent usage test
func TestConcurrentUsage(t *testing.T) {
	t.Run("FixedPool", func(t *testing.T) {
		pool := NewFixedPool(64)
		var wg sync.WaitGroup
		workers := 100
		iterations := 1000

		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()
				data := make([]byte, 32)

				for j := 0; j < iterations; j++ {
					buf := pool.Get()
					buf = append(buf, data...)
					pool.Put(buf)
				}
			}()
		}
		wg.Wait()
	})

	t.Run("GrowablePool", func(t *testing.T) {
		pool := NewGrowablePool()
		var wg sync.WaitGroup
		workers := 100
		iterations := 1000

		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()
				data := make([]byte, 32)

				for j := 0; j < iterations; j++ {
					buf := pool.Get()
					buf = append(buf, data...)
					pool.Put(buf)
				}
			}()
		}
		wg.Wait()
	})
}
