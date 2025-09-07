package main

import (
	"context"
	"os"

	"github.com/Nadim147c/rong/cmd"
	"github.com/charmbracelet/fang"
)

var Version = "dev"

func main() {
	if err := fang.Execute(context.Background(), cmd.Command, fang.WithVersion(Version)); err != nil {
		os.Exit(1)
	}
}
