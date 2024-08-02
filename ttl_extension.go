//go:build ttl
// +build ttl

package main

import (
	"net/http"
	"sync"
	"time"
)

type TTLExtension struct {
	ttl map[string]time.Time
	mu  sync.RWMutex
}

var ttlExtension = &TTLExtension{ttl: make(map[string]time.Time)}

func init() {
	extensions = append(extensions, ttlExtension)
}

func (e *TTLExtension) Init() {
	go e.cleanupExpiredKeys()
}

func (e *TTLExtension) BeforePut(key, value string, r *http.Request) {
	if ttlStr := r.URL.Query().Get("ttl"); ttlStr != "" {
		ttl, err := time.ParseDuration(ttlStr)
		if err == nil {
			e.mu.Lock()
			e.ttl[key] = time.Now().Add(ttl)
			e.mu.Unlock()
		}
	}
}

func (e *TTLExtension) AfterPut(key, value string, r *http.Request) {}

func (e *TTLExtension) BeforeGet(key string, r *http.Request) {
	e.mu.RLock()
	if expiry, exists := e.ttl[key]; exists && time.Now().After(expiry) {
		e.mu.RUnlock()
		e.deleteExpiredKey(key)
		return
	}
	e.mu.RUnlock()
}

func (e *TTLExtension) AfterGet(key string, value *string, r *http.Request) {}

func (e *TTLExtension) BeforeDelete(key string, r *http.Request) {}

func (e *TTLExtension) AfterDelete(key string, r *http.Request) {
	e.mu.Lock()
	delete(e.ttl, key)
	e.mu.Unlock()
}

func (e *TTLExtension) cleanupExpiredKeys() {
	for {
		time.Sleep(time.Minute)
		e.mu.Lock()
		now := time.Now()
		for key, expiry := range e.ttl {
			if now.After(expiry) {
				delete(e.ttl, key)
				store.mu.Lock()
				delete(store.data, key)
				store.mu.Unlock()
			}
		}
		e.mu.Unlock()
	}
}

func (e *TTLExtension) deleteExpiredKey(key string) {
	e.mu.Lock()
	delete(e.ttl, key)
	store.mu.Lock()
	delete(store.data, key)
	store.mu.Unlock()
	e.mu.Unlock()
}
