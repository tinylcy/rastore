package store

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewStore(t *testing.T) {
	s := NewStore()
	fmt.Println("s", s)
}

func TestStoreApply(t *testing.T) {
	var c1, c2 command
	var data []byte

	c1.Op = "set"
	c1.Key = "tinylcy"
	c1.Val = "chenyang"

	data, err := json.Marshal(c1)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal command: %s", err.Error()))
	}

	if err := json.Unmarshal(data, &c2); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal command: %s", err.Error()))
	}

	fmt.Println("command: ", c2)
}

func TestOpenStore(t *testing.T) {
	s := NewStore()
	s.RaftDir = "./node0"
	s.RaftBindAddr = ":11000"

	err := s.Open(true, "node0")
	if err != nil {
		panic(fmt.Sprintf("Faile to open store: %s", err.Error()))
	}
}
