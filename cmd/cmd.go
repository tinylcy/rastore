package cmd

import (
	"flag"
	"fmt"
	"os"
)

type Cmd struct {
	ServiceAddr string
	RaftAddr    string
	ClusterAddr string
	NodeID      string
	Data        string
}

func ParseCmd() *Cmd {
	var cmd Cmd
	flag.StringVar(&cmd.ServiceAddr, "serviceaddr", ":9090", "Set Rastore RESTFul service address.")
	flag.StringVar(&cmd.RaftAddr, "raftaddr", ":9091", "Set Raft communication address.")
	flag.StringVar(&cmd.ClusterAddr, "clusteraddr", "", "Set cluster address to be joined.")
	flag.StringVar(&cmd.NodeID, "id", "", "Set unique ID for node in cluster.")
	flag.StringVar(&cmd.Data, "data", ".", "Set Raft data path.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	return &cmd
}

func (cmd *Cmd) String() string {
	return fmt.Sprintf("[ServiceAddr -> %s, RaftAddr -> %s, ClusterAddr -> %s, NodeID -> %s, Data -> %s]",
		cmd.ServiceAddr,
		cmd.RaftAddr,
		cmd.ClusterAddr,
		cmd.NodeID,
		cmd.Data)
}
