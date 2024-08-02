//go:build persistent
// +build persistent

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

type PersistentExtension struct {
	filePath string
	mu       sync.RWMutex
}

var persistentExtension = &PersistentExtension{filePath: "store.json"}

func init() {
	extensions = append(extensions, persistentExtension)
}

func (e *PersistentExtension) Init() {
	e.load()
}

func (e *PersistentExtension) BeforePut(key, value string, r *http.Request) {}

func (e *PersistentExtension) AfterPut(key, value string, r *http.Request) {
	e.save()
}

func (e *PersistentExtension) BeforeGet(key string, r *http.Request) {}

func (e *PersistentExtension) AfterGet(key string, value *string, r *http.Request) {}

func (e *PersistentExtension) BeforeDelete(key string, r *http.Request) {}

func (e *PersistentExtension) AfterDelete(key string, r *http.Request) {
	e.save()
}

func (e *PersistentExtension) save() {
	e.mu.RLock()
	defer e.mu.RUnlock()

	store.mu.RLock()
	defer store.mu.RUnlock()

	file, err := os.Create(e.filePath)
	if err != nil {
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(store.data)
}

func (e *PersistentExtension) load() {
	e.mu.Lock()
	defer e.mu.Unlock()

	file, err := os.Open(e.filePath)
	if err != nil {
		return
	}
	defer file.Close()

	store.mu.Lock()
	defer store.mu.Unlock()

	json.NewDecoder(file).Decode(&store.data)
}
