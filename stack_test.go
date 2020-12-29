package lockfree

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mutexStack struct {
	mu sync.Mutex
	l  []int
}

func newMutextStack() *mutexStack {
	return &mutexStack{
		l: make([]int, 0),
	}
}

func (s *mutexStack) Push(item int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.l = append(s.l, item)
}

func (s *mutexStack) Pop() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := s.l[len(s.l)-1]
	s.l = s.l[:len(s.l)-1]

	return item
}

func TestMutexStack(t *testing.T) {
	s := newMutextStack()

	s.Push(1)
	s.Push(2)

	assert.Equal(t, s.Pop(), 2)
	assert.Equal(t, s.Pop(), 1)
}

func BenchmarkMutexStack(b *testing.B) {
	length := 1 << 12
	inputs := make([]int, length)

	for i := 0; i < length; i++ {
		inputs = append(inputs, rand.Int())
	}

	s := newMutextStack()
	b.ResetTimer()

	var cnt int64
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			i := int(atomic.AddInt64(&cnt, 1)-1) % length
			v := inputs[i]
			if v >= 0 {
				s.Push(v)
			} else {
				s.Pop()
			}
		}
	})
}

func BenchmarkLockfreeStack(b *testing.B) {
	length := 1 << 12
	inputs := make([]int, length)

	for i := 0; i < length; i++ {
		inputs = append(inputs, rand.Int())
	}

	s := NewStack()
	b.ResetTimer()

	var cnt int64
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			i := int(atomic.AddInt64(&cnt, 1)-1) % length
			v := inputs[i]
			if v >= 0 {
				s.Push(v)
			} else {
				s.Pop()
			}
		}
	})
}
