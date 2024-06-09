package service

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
)

type Options struct {
	Start StartOptions `command:"start" description:"start demon"`
	Stop  StopOptions  `command:"stop" description:"stop the demon"`
}

type StartOptions struct {
	Yes bool `short:"y" long:"yes" description:"yes to all"`
}
type StopOptions struct{}

func ApplyCommand(cmd *flag.Command, data Options) {
	if cmd.Active.Name != "service" {
		return
	}

	config.WithCommit(func() bool {
		switch cmd.Active.Active.Name {
		case "start":
			return StartService(data.Start.Yes)
		case "stop":
			return StopService()
		default:
			return false
		}
	})
}
