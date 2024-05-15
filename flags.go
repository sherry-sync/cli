package main

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/auth"
	"sherry/shr/folder"
)

type Options struct {
	ConfigPath flag.Filename  `long:"config" short:"c" description:"Path to configuration folder"`
	Auth       auth.Options   `command:"auth" description:"Authenticate"`
	Folder     folder.Options `command:"folder" description:"Folder operations"`
}

func applyCommand(cmd *flag.Command, options Options) {
	auth.ApplyCommand(cmd, options.Auth)
	folder.ApplyCommands(cmd, options.Folder)
}
