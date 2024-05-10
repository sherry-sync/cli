package folder

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
	"sherry/shr/helpers"
)

var CommandGroupName = "folder"
var Flags struct {
	Path flag.Filename     `long:"path" short:"p" description:"Specify local path for operation"`
	User string            `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Set  map[string]string `long:"set" short:"s" description:"Set folder settings"`

	Create bool   `long:"create" short:"c" description:"Create shared folder"`
	Name   string `long:"name" short:"n" description:"Specify shared folder name"`
	Yes    bool   `long:"yes" short:"y" description:"Skip confirmation and use default values where possible"`

	Get bool `long:"get" short:"g" description:"Get shared folder"`

	Display bool `long:"display" short:"d" description:"Display shared folder shearing info"`
}

func InitCommands(cmd *flag.Command) {
	helpers.AddCommand(cmd,
		CommandGroupName,
		"Folder",
		"Manage your folders",
		&Flags,
	)
}

func ApplyCommands(cmd *flag.Command) {
	if cmd.Active.Name != CommandGroupName {
		return
	}
	helpers.PrintJson(Flags)
	helpers.PrintJson(config.GetConfig())
}
