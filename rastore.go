package main

import (
	"fmt"
	"log"
	"os"

	"rastore/cmd"
	"rastore/service"
	"rastore/store"
)

func main() {
	cmd := cmd.ParseCmd()
	if cmd.Data == "" {
		fmt.Fprintf(os.Stderr, "Please specify the Raft data path.")
		return
	}
	if err := os.Mkdir(cmd.Data, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Faile to create directory for Raft data storage: %s", err.Error())
		return
	}

	// Create Raft storer and set bind address & data storage path.
	st := store.NewStore()
	st.RaftDir = cmd.Data
	st.RaftBindAddr = cmd.RaftAddr

	// Start Raft consensus.
	if err := st.Open(cmd.ClusterAddr == "", cmd.NodeID); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open the Raft consensus storer: %s", err.Error())
		return
	}

	// Create key-value service provider.
	se := service.NewService(cmd.ServiceAddr, st)
	if err := se.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open key-value service provider: %s", err.Error())
		return
	}

	log.Print("Rastore started successfully.")

}
