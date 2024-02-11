package domain

import (
	"sync"
	"sync/atomic"
)

type AtomicSignatureCounter struct {
	counter int32
}

var (
	counterInstance *AtomicSignatureCounter
	once            sync.Once
)

func NewAtomicSignatureCounter() *AtomicSignatureCounter {
	return &AtomicSignatureCounter{counter: -1}
}

func Increment() *AtomicSignatureCounter {
	once.Do(func() {
		counterInstance = NewAtomicSignatureCounter()
	})
	atomic.AddInt32(&counterInstance.counter, 1)
	return counterInstance
}

func (c *AtomicSignatureCounter) Get() int32 {
	return atomic.LoadInt32(&c.counter)
}
