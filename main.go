package main

import (
	"fmt"
	"os"

	"github.com/Nadim147c/rong/cmd"
)

func main() {
	if err := cmd.Command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
