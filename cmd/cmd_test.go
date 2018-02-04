package cmd

import (
	"fmt"
	"testing"
)

func TestDefaultParseCmd(t *testing.T) {
	cmd := ParseCmd()
	fmt.Println(cmd)
}
