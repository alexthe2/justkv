package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlePut(t *testing.T) {
	req := httptest.NewRequest("PUT", "/put/testkey", strings.NewReader("testvalue"))
	w := httptest.NewRecorder()
	handlePut(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %v", res.Status)
	}

	store.mu.RLock()
	defer store.mu.RUnlock()
	if store.data["testkey"] != "testvalue" {
		t.Errorf("expected value 'testvalue', got '%v'", store.data["testkey"])
	}
}

func TestHandleGet(t *testing.T) {
	store.mu.Lock()
	store.data["testkey"] = "testvalue"
	store.mu.Unlock()

	req := httptest.NewRequest("GET", "/get/testkey", nil)
	w := httptest.NewRecorder()
	handleGet(w, req)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %v", res.Status)
	}

	if string(body) != "testvalue" {
		t.Errorf("expected value 'testvalue', got '%v'", string(body))
	}
}

func TestHandleDelete(t *testing.T) {
	store.mu.Lock()
	store.data["testkey"] = "testvalue"
	store.mu.Unlock()

	req := httptest.NewRequest("DELETE", "/delete/testkey", nil)
	w := httptest.NewRecorder()
	handleDelete(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %v", res.Status)
	}

	store.mu.RLock()
	defer store.mu.RUnlock()
	if _, exists := store.data["testkey"]; exists {
		t.Errorf("expected key 'testkey' to be deleted")
	}
}
