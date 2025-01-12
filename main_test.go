package main

import "testing"

func BenchmarkSliceAllocateCapacity(b *testing.B) {
	b.ResetTimer()
	const size int16 = 1000
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		SliceAllocateCapacity(size)
	}
}

func BenchmarkSliceLetCapacityGrow(b *testing.B) {
	b.ResetTimer()
	const size int16 = 1000
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		SliceLetCapacityGrow(size)
	}
}
