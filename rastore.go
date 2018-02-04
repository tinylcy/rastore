package main

import (
	"fmt"

	"rastore/cmd"
)

func main() {
	cmd := cmd.ParseCmd()
	fmt.Println(cmd)
}
