package pkg

import (
	"sync"
)

type Slice[T any] struct {
	mu    sync.Mutex
	slice []T
}

func (s *Slice[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.slice)
}

func (s *Slice[T]) Append(i T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.slice = append(s.slice, i)
}

func (s *Slice[T]) Foreach(f func(n int, i T)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for n, i := range s.slice {
		f(n, i)
	}
}
