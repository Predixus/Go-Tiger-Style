package main

func SliceAllocateCapacity(n int16) {
	data := make([]int16, 0, n)
	for i := 0; i < int(n); i++ {
		data = append(data, int16(i)) // No reallocation needed
	}
}

func SliceLetCapacityGrow(n int16) {
	data := make([]int16, 0)
	for i := 0; i < int(n); i++ {
		data = append(data, int16(i)) // May trigger reallocation
	}
}
