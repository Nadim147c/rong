package main

import (
	"github.com/Nadim147c/rong/cmd"
	cc "github.com/ivanpirog/coloredcobra"
)

func main() {
	cc.Init(&cc.Config{
		RootCmd:         cmd.Command,
		Headings:        cc.Cyan + cc.Bold + cc.Underline,
		Commands:        cc.Yellow + cc.Bold,
		CmdShortDescr:   cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Green + cc.Bold,
		FlagsDataType:   cc.Red + cc.Bold,
		NoExtraNewlines: true,
	})

	cmd.Command.Execute()
}
