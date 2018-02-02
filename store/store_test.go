package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// Test that the store can be opened.
func TestStoreOpen(t *testing.T) {
	s := New()
	tempDir, _ := ioutil.TempDir("", "store_test")
	defer os.Remove(tempDir)

	s.RaftBind = "127.0.0.1:0"
	s.RaftDir = tempDir

	if s == nil {
		t.Fatalf("failed to create store")
	}

	if err := s.Open(false, "node0"); err != nil {
		t.Fatalf("failed to open store: %s", err)
	}
}

func TestStoreOpenSingleNode(t *testing.T) {
	s := New()
	tempDir, _ := ioutil.TempDir("", "store_test")
	defer os.RemoveAll(tempDir)

	s.RaftBind = "127.0.0.1:0"
	s.RaftDir = tempDir
	if s == nil {
		t.Fatalf("failed to create store")
	}

	if err := s.Open(true, "node0"); err != nil {
		t.Fatalf("faile to open store: %s", err)
	}

	// Simple way to ensure there is a leader.
	time.Sleep(3 * time.Second)

	if err := s.Set("foo", "bar"); err != nil {
		t.Fatalf("failed to set key: %s", err.Error())
	}

	// Wait for committed log entry to be applied.
	time.Sleep(500 * time.Millisecond)

	value, err := s.Get("foo")
	if err != nil {
		t.Fatalf("faile to get key: %s", err.Error())
	}

	if value != "bar" {
		t.Fatalf("key has wrong value: %s", value)
	}

	if err := s.Delete("foo"); err != nil {
		t.Fatalf("failed to delete key: %s", err.Error())
	}

	// Wait for committed log entry to be applied.
	time.Sleep(500 * time.Millisecond)
	value, err = s.Get("foo")
	if err != nil {
		t.Fatalf("faile to get key: %s", err.Error())
	}

	if value != "" {
		t.Fatalf("key has wrong value: %s", value)
	}
}
