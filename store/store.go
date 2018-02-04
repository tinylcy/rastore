package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

const (
	raftTimeOut = 10 * time.Second // The timeout to limit the amount of time we wait for the command to be started.
)

type Store struct {
	RaftDir      string
	RaftBindAddr string

	data map[string]string

	raft *raft.Raft // Raft consensus algorithm.

	mu sync.Mutex
}

type command struct {
	Op  string
	Key string
	Val string
}

// Create a new store without Raft initialization.
func NewStore() *Store {
	return &Store{
		RaftDir:      "",
		RaftBindAddr: "",
		data:         make(map[string]string),
	}
}

// Open the store.
// If enableSingle is set and there are no existing peers, the node will
// become the first node and therefore become the leader of the cluster.
func (s *Store) Open(enableSingle bool, localID string) error {
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID) // LocalID is the unique ID for this node across all time. API: type ServerID string.

	// Setup Raft network communication.
	// Create a NetworkTransport that is build in on top of a TCP streaming transport layer.
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBindAddr)
	if err != nil {
		return errors.New(fmt.Sprintf("Faile to resolve TCP address: %s", err.Error()))
	}
	transport, err := raft.NewTCPTransport(s.RaftBindAddr, addr, 3, raftTimeOut, os.Stderr)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create TCP transport: %s", err.Error()))
	}

	// Create the snapshot store, which allows the Raft to truncate the log.
	snapshot, err := raft.NewFileSnapshotStore(s.RaftDir, 3, os.Stderr)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create file snapshot store: %s", err.Error()))
	}

	// Create the log store and stable store.
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft.db"))
	if err != nil {
		return errors.New(fmt.Sprintf("Faied to create log store and stable store: %s", err.Error()))
	}

	// Initialize the Raft system.
	r, err := raft.NewRaft(config, s, logStore, logStore, snapshot, transport)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create Raft system: %s", err.Error()))
	}
	s.raft = r

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}

		// This should only be called at the beginning of time for the cluster.
		r.BootstrapCluster(configuration)
	}

	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[key], nil
}

func (s *Store) Set(key string, val string) error {
	// Check current raft node state.
	// All changes to the system must go through the leader.
	if s.raft.State() != raft.Leader {
		return errors.New("Not the leader, all changes to the system must go through the leader.")
	}

	c := &command{
		Op:  "SET",
		Key: key,
		Val: val,
	}

	cmd, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Apply is used to issue a command to the FSM in a highly consistent manner.
	// This returns a future that ca be used to wait on the application.
	// This must be run on the leader or it will fail.
	future := s.raft.Apply(cmd, raftTimeOut)

	return future.Error()
}

func (s *Store) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return errors.New("Not the leader, all changes to the system must go through the leader.")
	}

	c := &command{
		Op:  "DELETE",
		Key: key,
	}

	cmd, err := json.Marshal(c)
	if err != nil {
		return err
	}

	future := s.raft.Apply(cmd, raftTimeOut)

	return future.Error()
}

////////////////////////////// FSM /////////////////////////////////
// FSM provides an interface that can be implementated by clients //
// to make use of replicated log.                                 //
////////////////////////////// FSM /////////////////////////////////

// Apply log is invoked once a log enetry is committed.
func (s *Store) Apply(log *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(log.Data, &c); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case "SET":
		return s.applySet(c.Key, c.Val)
	case "DELETE":
		return s.applyDelete(c.Key)
	default:
		return fmt.Sprintf("Unsupported operation: %s", c.Op)
	}

	return nil
}

// Snapshot is used to support log compaction. This call should return
// an FSMSnapshot which can be used to save a point-in-time snapshot of the FSM.
func (s *Store) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

// Restore is used to restore an FSM from a snapshot.
func (s *Store) Restore(rc io.ReadCloser) error {
	return nil
}

func (s *Store) applySet(key string, val string) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
	return nil
}

func (s *Store) applyDelete(key string) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}
