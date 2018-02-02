package store

import (
	"io/ioutil"
	"os"
	"testing"
)

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
