package cache

import (
	"math/rand"
	"sync"
)

type Cache[V any] struct {
	sl []V
	mu *sync.RWMutex
}

func New[V any](size int) *Cache[V] {
	return &Cache[V]{
		sl: make([]V, 0, size),
		mu: &sync.RWMutex{},
	}
}

func (c *Cache[V]) Fill(data []V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.sl = append(c.sl, data...)
}

func (c *Cache[V]) GetRandom() V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var v V
	if len(c.sl) == 0 {
		return v
	}

	return c.sl[rand.Intn(len(c.sl)-1)]
}
