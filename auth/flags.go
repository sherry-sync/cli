package auth

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
	"sherry/shr/helpers"
)

var CommandGroupName = "auth"
var Flags struct {
	Register bool   `long:"register" short:"r" description:"Register new user"`
	List     bool   `long:"list" short:"l" description:"List authorized users"`
	User     string `long:"username" short:"u" description:"Username"`
	Email    string `long:"email" short:"e" description:"User email"`
	Password string `long:"password" short:"p" description:"User password"`
	Default  string `long:"default" short:"d" description:"Set default user or get (--default=)" default:"-" default-mask:"-"`
}

func InitCommand(cmd *flag.Command) {
	helpers.AddCommand(cmd,
		CommandGroupName,
		"Authorization",
		"Manage your authorization",
		&Flags,
	)
}

func ApplyCommand(cmd *flag.Command) {
	if cmd.Active.Name != CommandGroupName {
		return
	}

	shouldCommit := false
	if Flags.Register {
		shouldCommit = RegisterUser(Flags.Email, Flags.Password, Flags.User)
	} else if Flags.Default != "-" {
		if Flags.Default == "" {
			shouldCommit = PrintDefaultUser()
		} else {
			shouldCommit = SetDefaultUser(Flags.Default)
		}
	} else if Flags.List {
		shouldCommit = PrintUsers()
	} else {
		shouldCommit = LoginUser(Flags.Email, Flags.Password)
	}

	if shouldCommit {
		config.CommitAuth()
	}
}
