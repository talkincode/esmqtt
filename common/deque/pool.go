package deque

import (
	"sync"
)

type Pool[T any] struct {
	quemap map[string]*Deque[T]
	sync.Mutex
}

func NewPool[T any]() *Pool[T] {
	return &Pool[T]{
		quemap: make(map[string]*Deque[T]),
	}
}

// GetQueue get queue by name
func (p *Pool[T]) GetQueue(name string, size int) *Deque[T] {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.quemap[name]; !ok {
		p.quemap[name] = New[T](size)
	}
	return p.quemap[name]
}

// RemoveQueue remove queue by name
func (p *Pool[T]) RemoveQueue(name string) {
	p.Lock()
	defer p.Unlock()
	delete(p.quemap, name)
}

// ExistQueue exist queue by name
func (p *Pool[T]) ExistQueue(name string) bool {
	p.Lock()
	defer p.Unlock()
	_, ok := p.quemap[name]
	return ok
}
