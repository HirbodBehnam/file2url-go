package database

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

// MemoryCache is an in memory cache to store files which want to be downloaded temporary
// The cache is deleted every hour
type MemoryCache struct {
	m  map[string]File
	mu sync.RWMutex
}

// NewMemoryCache creates a new memory cache and setups a cleanup goroutine
func NewMemoryCache() *MemoryCache {
	m := &MemoryCache{m: make(map[string]File)}
	go m.cleanupGoroutine()
	return m
}

// cleanupGoroutine deletes old files from memory
func (m *MemoryCache) cleanupGoroutine() {
	for {
		time.Sleep(time.Hour)
		m.mu.Lock()
		start := time.Now()
		for k, v := range m.m {
			// Delete entries older than a day
			if start.Sub(v.AddedTime) > time.Hour*24 {
				delete(m.m, k)
			}
		}
		m.mu.Unlock()
	}
}

// Store stores the file in the cache and returns an ID for it
func (m *MemoryCache) Store(f File) (string, error) {
	id := uuid.NewString()
	m.mu.Lock()
	m.m[id] = f
	m.mu.Unlock()
	return id, nil
}

// Load loads a file from cache
func (m *MemoryCache) Load(id string) (File, bool) {
	m.mu.RLock()
	f, exists := m.m[id]
	m.mu.RUnlock()
	return f, exists
}

// Close is no-op in here
func (m *MemoryCache) Close() error {
	return nil
}
