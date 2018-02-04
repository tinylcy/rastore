package cmd

import (
	"flag"
	"fmt"
	"os"
)

type Cmd struct {
	serviceAddr string
	raftAddr    string
	clusterAddr string
	nodeID      string
	data        string
}

func ParseCmd() *Cmd {
	var cmd Cmd
	flag.StringVar(&cmd.serviceAddr, "serviceaddr", ":9090", "Set Rastore RESTFul service address.")
	flag.StringVar(&cmd.raftAddr, "raftaddr", ":9091", "Set Raft communication address.")
	flag.StringVar(&cmd.clusterAddr, "clusteraddr", "", "Set cluster address to be joined.")
	flag.StringVar(&cmd.nodeID, "id", "", "Set unique ID for node in cluster.")
	flag.StringVar(&cmd.data, "data", ".", "Set Raft data path.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	return &cmd
}

func (cmd *Cmd) String() string {
	return fmt.Sprintf("[ServiceAddr -> %s, RaftAddr -> %s, ClusterAddr -> %s, NodeID -> %s, Data -> %s]",
		cmd.serviceAddr,
		cmd.raftAddr,
		cmd.clusterAddr,
		cmd.nodeID,
		cmd.data)
}
