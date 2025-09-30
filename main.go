package main

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/Nadim147c/rong/cmd"
	"github.com/charmbracelet/fang"
)

var Version = "dev"

func main() {
	err := fang.Execute(
		context.Background(),
		cmd.Command,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme),
		fang.WithNotifySignal(os.Interrupt, os.Kill),
		fang.WithVersion(Version),
		fang.WithoutCompletions(),
		fang.WithErrorHandler(func(w io.Writer, styles fang.Styles, err error) {
			if errors.Is(err, context.Canceled) {
				err = errors.New("operation cancelled by user")
			}
			fang.DefaultErrorHandler(w, styles, err)
		}),
	)
	if err != nil {
		os.Exit(1)
	}
}
