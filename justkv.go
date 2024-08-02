package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// Store to hold the key-value pairs
type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

var store = Store{data: make(map[string]string)}

// Extension interface
type Extension interface {
	Init()
	BeforePut(key, value string, r *http.Request)
	AfterPut(key, value string, r *http.Request)
	BeforeGet(key string, r *http.Request)
	AfterGet(key string, value *string, r *http.Request)
	BeforeDelete(key string, r *http.Request)
	AfterDelete(key string, r *http.Request)
}

var extensions []Extension

func main() {
	// Initialize extensions
	for _, ext := range extensions {
		ext.Init()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/put/", handlePut)
	http.HandleFunc("/get/", handleGet)
	http.HandleFunc("/delete/", handleDelete)

	fmt.Println("Server started at:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/put/"):]
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, ext := range extensions {
		ext.BeforePut(key, string(value), r)
	}

	store.mu.Lock()
	store.data[key] = string(value)
	store.mu.Unlock()

	for _, ext := range extensions {
		ext.AfterPut(key, string(value), r)
	}

	w.WriteHeader(http.StatusOK)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/get/"):]
	for _, ext := range extensions {
		ext.BeforeGet(key, r)
	}

	store.mu.RLock()
	value, exists := store.data[key]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	for _, ext := range extensions {
		ext.AfterGet(key, &value, r)
	}

	w.Write([]byte(value))
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/delete/"):]
	for _, ext := range extensions {
		ext.BeforeDelete(key, r)
	}

	store.mu.Lock()
	_, exists := store.data[key]
	if exists {
		delete(store.data, key)
	}
	store.mu.Unlock()

	for _, ext := range extensions {
		ext.AfterDelete(key, r)
	}

	if exists {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Key not found", http.StatusNotFound)
	}
}
