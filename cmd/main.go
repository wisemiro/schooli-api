package main

import (
	"os"
	"schooli-api/cmd/commands"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}

// TODO: Resource consumption
