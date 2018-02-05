package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"rastore/cmd"
	"rastore/service"
	"rastore/store"
)

func main() {
	cmd := cmd.ParseCmd()
	if cmd.Data == "" {
		fmt.Fprintf(os.Stderr, "Please specify the Raft data storage path\n")
		return
	}
	if err := os.Mkdir(cmd.Data, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Faile to create directory for Raft data storage: %s\n", err.Error())
		return
	}

	// Create Raft storer and set bind address & data storage path.
	st := store.NewStore()
	st.RaftDir = cmd.Data
	st.RaftBindAddr = cmd.RaftAddr

	// Start Raft consensus.
	if err := st.Open(cmd.ClusterAddr == "", cmd.NodeID); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open the Raft consensus storer: %s\n", err.Error())
		return
	}

	// Create key-value service provider.
	se := service.NewService(cmd.ServiceAddr, st)
	if err := se.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open key-value service provider: %s\n", err.Error())
		return
	}

	// Join cluster.
	if cmd.ClusterAddr != "" {
		join(cmd.NodeID, cmd.RaftAddr, cmd.ClusterAddr)
	}

	// Set up channel on which to send signal notification.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Block until a signal is received.
	<-interrupt
	fmt.Fprintf(os.Stderr, "\nExit.\n")
}

func join(nodeID string, nodeAddr string, clusterAddr string) {
	var url = fmt.Sprintf("http://%s/rastore/join", clusterAddr)
	log.Printf("Join a cluster: %s\n", url)

	nodes := map[string]string{nodeID: nodeAddr}
	jsonNodes, err := json.Marshal(nodes)
	if err != nil {
		log.Fatalf("Faile to marshal nodes value: %s", err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonNodes))
	if err != nil {
		log.Fatalf("Faild to send POST request to join cluster: %s", err.Error())
	}
	defer resp.Body.Close()
}
