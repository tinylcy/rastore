package service

import (
	"fmt"
	"testing"
)

type testStore struct {
	data map[string]string
}

func NewTestStore() *testStore {
	return &testStore{
		data: make(map[string]string),
	}
}

func (t *testStore) Get(key string) (string, error) {
	return t.data[key], nil
}

func (t *testStore) Set(key string, val string) error {
	t.data[key] = val
	return nil
}

func (t *testStore) Delete(key string) error {
	delete(t.data, key)
	return nil
}

func TestOpenService(t *testing.T) {
	store := NewTestStore()
	listenAddr := ":8080"
	service := NewService(listenAddr, store)
	fmt.Println("service:", service)
	service.Open()
}
