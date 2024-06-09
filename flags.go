package main

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/auth"
	"sherry/shr/folder"
	"sherry/shr/service"
)

type Options struct {
	ConfigPath flag.Filename   `long:"config" short:"c" description:"Path to configuration folder"`
	Auth       auth.Options    `command:"auth" description:"Authenticate"`
	Folder     folder.Options  `command:"folder" description:"Folder operations"`
	Service    service.Options `command:"service" description:"service operations"`
}

func applyCommand(cmd *flag.Command, options Options) {
	auth.ApplyCommand(cmd, options.Auth)
	folder.ApplyCommands(cmd, options.Folder)
	service.ApplyCommand(cmd, options.Service)
}
