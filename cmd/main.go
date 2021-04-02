package cmd

import (
	"fmt"
	"os"
)

func Run() {

	args := os.Args

	if len(args) == 1 {
		fmt.Fprintln(os.Stderr, "Please provide a valid subcommand")
		//Print usage information
		os.Exit(1)
	}

	command := args[1]

	switch command {
	case "version":
		fmt.Println("Development version")
	case "convert":
		convertCmd := newConvertCommand()
		convertCmd.run(os.Args[2:])
	}
}
