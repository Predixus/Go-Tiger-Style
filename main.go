package main

import (
	"fmt"
	"sync"
)

func SliceAllocateCapacity() {
	data := make([]int16, 0, 1000)
	for i := 0; i < 1000; i++ {
		data = append(data, int16(i)) // No reallocation needed
	}
}

func SliceLetCapacityGrow() {
	data := make([]int16, 0)
	for i := 0; i < 1000; i++ {
		data = append(data, int16(i)) // May trigger reallocation
	}
}

func SliceNoShareCapacity() []int16 {
	original := make([]int16, 5, 5)
	slice2 := make([]int16, 3)
	copy(slice2, original[0:3]) // Copy just the elements we want
	slice2 = append(slice2, 6)  // Now this truly won't affect original
	return original
}

func SliceShareCapacity() []int16 {
	original := make([]int16, 5, 10)
	slice2 := original[0:3]    // slice2 has capacity of 7
	slice2 = append(slice2, 6) // Modifies original's backing array!
	return original
}

type User struct{}

func MapSingleHashAllocation() {
	users := make(map[string]User, 1000)
	for i := 0; i < 1000; i++ {
		users[fmt.Sprintf("user%d", i)] = User{} // No rehashing needed
	}
}

func MapMultipleRehashings() {
	users := make(map[string]User) // Default small capacity
	for i := 0; i < 1000; i++ {
		users[fmt.Sprintf("user%d", i)] = User{} // Forces periodic rehashing
	}
}

// Channel operations demonstrating sync vs async communication
func ChannelSync() int {
	chSync := make(chan int)
	go func() {
		chSync <- 1 // Blocks until receiver is ready
	}()
	return <-chSync // Blocks until sender sends
}

func ChannelAsync() int {
	chAsync := make(chan int, 1)
	chAsync <- 1     // Doesn't block because buffer available
	return <-chAsync // Dequeue from buffer
}

// Multiple operations to demonstrate blocking behavior
func ChannelSyncMulti() []int {
	chSync := make(chan int)
	done := make(chan bool)
	result := make([]int, 0, 100) // Pre-allocate capacity

	go func() {
		for i := 0; i < 100; i++ {
			chSync <- i // Each send blocks until received
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			result = append(result, <-chSync) // Each receive blocks until sent
		}
		done <- true
	}()

	<-done
	<-done
	return result
}

func ChannelAsyncMulti() []int {
	chAsync := make(chan int, 100)
	done := make(chan bool)
	result := make([]int, 0, 100) // Pre-allocate capacity

	go func() {
		for i := 0; i < 100; i++ {
			chAsync <- i // Won't block until buffer full
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			result = append(result, <-chAsync) // Won't block while buffer has items
		}
		done <- true
	}()

	<-done
	<-done
	return result
}

type FixedPool struct {
	buffers sync.Pool
	size    int
}

func NewFixedPool(size int) *FixedPool {
	return &FixedPool{
		buffers: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, size)
			},
		},
		size: size,
	}
}

func (p *FixedPool) Get() []byte {
	buf := p.buffers.Get().([]byte)
	return buf[:0] // Reset length but keep capacity
}

func (p *FixedPool) Put(buf []byte) {
	if cap(buf) == p.size {
		p.buffers.Put(buf)
	}
	// Discard buffers that have grown beyond our fixed size
}

// GrowablePool allows buffers to grow
type GrowablePool struct {
	buffers sync.Pool
}

func NewGrowablePool() *GrowablePool {
	return &GrowablePool{
		buffers: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0)
			},
		},
	}
}

func (p *GrowablePool) Get() []byte {
	return p.buffers.Get().([]byte)[:0]
}

func (p *GrowablePool) Put(buf []byte) {
	p.buffers.Put(buf)
}
