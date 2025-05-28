package cmd

import (
	"github.com/Nadim147c/rong/cmd/image"
	"github.com/spf13/cobra"
)

// Command is root command of the cli
var Command = &cobra.Command{
	Use:   "rong",
	Short: "A material you color generator from image.",
}

func init() {
	Command.Flags()
	Command.AddCommand(image.Command)
}
