package folder

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
	"sherry/shr/helpers"
)

var CommandGroupName = "folder"
var Flags struct {
	User   string            `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Path   flag.Filename     `long:"path" short:"p" description:"Specify local path for operation"`
	Create bool              `long:"create" short:"c" description:"Create shared folder"`
	Get    bool              `long:"get" short:"g" description:"Get shared folder"`
	Set    map[string]string `long:"set" short:"s" description:"Set folder settings"`
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
