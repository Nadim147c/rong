package main

import (
	"context"
	"os"

	"github.com/Nadim147c/rong/cmd"
	"github.com/charmbracelet/fang"
)

func main() {
	if err := fang.Execute(context.Background(), cmd.Command); err != nil {
		os.Exit(1)
	}
}
