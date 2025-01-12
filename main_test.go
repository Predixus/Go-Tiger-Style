package main

import (
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
