package auth

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/cli/config"
	"sherry/cli/helpers"
)

var CommandGroupName = "auth"
var Flags struct {
	Register bool   `long:"register" short:"r" description:"Register new user"`
	User     string `long:"user" short:"u" description:"Username"`
	Password string `long:"password" short:"p" description:"User password"`
	Default  string `long:"default" short:"d" description:"Get/Set default user" default:"-" default-mask:"-"`
}

func InitCommand(cmd *flag.Command) {
	helpers.AddCommand(cmd,
		CommandGroupName,
		"Authorization",
		"Manage your authorization",
		&Flags,
	)
}

func ApplyCommand(cmd *flag.Command, config config.Config) {
	if cmd.Active.Name != CommandGroupName {
		return
	}

	if Flags.Register {
		RegisterUser(Flags.User, Flags.Password)
		return
	}
	if Flags.Default != "-" {
		SetDefaultUser(Flags.Default)
		return
	}

	LoginUser(Flags.User, Flags.Password)
}
