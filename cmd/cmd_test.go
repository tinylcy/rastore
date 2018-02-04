package cmd

import (
	"fmt"
	"testing"
)

func TestDefaultParseCmd(t *testing.T) {
	cmd := ParseCmd()
	fmt.Println("serviceAddr: ", cmd.serviceAddr)
	fmt.Println("raftAddr: ", cmd.raftAddr)
	fmt.Println("clusterAddr: ", cmd.clusterAddr)
	fmt.Println("nodeID: ", cmd.nodeID)
	fmt.Println("data: ", cmd.data)
}
