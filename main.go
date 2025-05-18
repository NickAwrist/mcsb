package main

import (
	"fmt"
	"mcsb-cli/cmd"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "create" {
			cmd.Execute()
		}
	} else {
		fmt.Println("Usage: mcsb create")
	}
}
