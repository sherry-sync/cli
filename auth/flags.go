package auth

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
)

type Options struct {
	Register RegisterOptions `command:"register" description:"Register new user"`
	Login    LoginOptions    `command:"login" description:"Authorize existing user"`
	Default  DefaultOptions  `command:"default" description:"Display/Set default user"`
	List     List            `command:"list" description:"List authorized users"`
}

type RegisterOptions struct {
	User     string `long:"username" short:"u" description:"Username"`
	Email    string `long:"email" short:"e" description:"User email"`
	Password string `long:"password" short:"p" description:"User password"`
}

type LoginOptions struct {
	Email    string `long:"email" short:"e" description:"User email"`
	Password string `long:"password" short:"p" description:"User password"`
}

type DefaultOptions struct {
	Args struct {
		Username string `positional-arg-name:"username" description:"Optional username to set as default user"`
	} `positional-args:"yes" description:"Optional username to set as default user"`
}

type List struct {
}

func ApplyCommand(cmd *flag.Command, data Options) {
	if cmd.Active.Name != "auth" {
		return
	}

	config.WithCommit(func() bool {
		switch cmd.Active.Active.Name {
		case "register":
			return RegisterUser(data.Register.Email, data.Register.Password, data.Register.User)
		case "login":
			return LoginUser(data.Login.Email, data.Login.Password)
		case "list":
			return PrintUsers()
		case "default":
			var username = data.Default.Args.Username
			if username == "" {
				return PrintDefaultUser()
			} else {
				return SetDefaultUser(username)
			}
		default:
			return false
		}
	})
}
