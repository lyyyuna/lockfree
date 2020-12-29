package lockfree

import (
	"math"
	"sync/atomic"
	"unsafe"
)

type Stack struct {
	top unsafe.Pointer
	len uint64
}

type stackNode struct {
	val  int
	next unsafe.Pointer
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(v int) {
	item := stackNode{
		val: v,
	}
	var top unsafe.Pointer
	for {
		top = atomic.LoadPointer(&s.top)
		item.next = top
		if atomic.CompareAndSwapPointer(&s.top, s.top, unsafe.Pointer(&item)) {
			atomic.AddUint64(&s.len, 1)
			return
		}
	}
}

func (s *Stack) Pop() int {
	var top, next unsafe.Pointer
	var item *stackNode
	for {
		top = atomic.LoadPointer(&s.top)
		// stack is empty
		if top == nil {
			return math.MaxInt32
		}

		item = (*stackNode)(top)
		next = atomic.LoadPointer(&item.next)

		if atomic.CompareAndSwapPointer(&s.top, s.top, next) {
			atomic.AddUint64(&s.len, ^uint64(0))
			return item.val
		}
	}
}
