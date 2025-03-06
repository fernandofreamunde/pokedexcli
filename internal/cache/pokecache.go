package internal

import (
	"fmt"
	"sync"
	"time"
)

type PokeCache struct {
	entries map[string]CacheEntry
	mu      sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval int) *PokeCache {
	c := PokeCache{}
	c.entries = map[string]CacheEntry{}

	go c.reapLoop(interval)

	return &c
}

func (c *PokeCache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("Cache added for: " + key)

	c.entries[key] = CacheEntry{val: val, createdAt: time.Now()}
}

func (c *PokeCache) Get(key string) ([]byte, bool) {

	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.entries[key]
	if !ok {
		fmt.Println("Cache miss!")
		return []byte{}, false
	}

	fmt.Println("Cache hit!")

	return value.val, true
}

func (c *PokeCache) reapLoop(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for i, entry := range c.entries {
				if now.After(entry.createdAt.Add(time.Duration(interval) * time.Second)) {
					delete(c.entries, i)
				}
			}
			c.mu.Unlock()
		}
	}
}
