package main

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/auth"
	"sherry/shr/folder"
)

var defaultFlags struct {
	ConfigPath flag.Filename `long:"config" short:"c" description:"Path to configuration folder"`
}

func initCommand(cmd *flag.Command) {
	auth.InitCommand(cmd)
	folder.InitCommands(cmd)
}

func applyCommand(cmd *flag.Command) {
	auth.ApplyCommand(cmd)
	folder.ApplyCommands(cmd)
}
