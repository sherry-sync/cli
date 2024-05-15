package folder

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
)

type Options struct {
	Create CreateOptions `command:"create" description:"Create shared folder"`
	Get    GetOptions    `command:"get" description:"Get shared folder"`
	Show   ShowOptions   `command:"show" description:"Display shared folder sharing info"`
	Set    SetOptions    `command:"set" description:"Set shared folder settings"`
}

type CreateOptions struct {
	Path flag.Filename     `long:"path" short:"p" description:"Specify local path for operation"`
	User string            `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Name string            `long:"name" short:"n" description:"Specify shared folder name"`
	Yes  bool              `long:"yes" short:"y" description:"Skip confirmation and use default values where possible"`
	Set  map[string]string `long:"set" short:"s" description:"Set folder settings"`
}

type GetOptions struct {
	Path flag.Filename `long:"path" short:"p" description:"Specify local path for operation"`
	User string        `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Name string        `long:"name" short:"n" description:"Shared folder in format owner_username:folder_name or folder id"`
	Yes  bool          `long:"yes" short:"y" description:"Skip confirmation and use default values where possible"`
}

type ShowOptions struct {
	User string `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Name string `long:"name" short:"n" description:"Shared folder in format owner_username:folder_name or folder id"`
}

type SetOptions struct {
	User string            `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Name string            `long:"name" short:"n" description:"Shared folder in format owner_username:folder_name or folder id"`
	Set  map[string]string `positional-arg-name:"set"  description:"Set folder settings"`
}

func ApplyCommands(cmd *flag.Command, options Options) {
	if cmd.Active.Name != "folder" {
		return
	}

	config.WithCommit(func() bool {
		switch cmd.Active.Active.Name {
		case "create":
			return CreateSharedFolder(options.Create.User, options.Create.Yes, string(options.Create.Path), options.Create.Name, options.Create.Set)
		case "get":
			return GetSharedFolder(options.Get.User, options.Get.Yes, string(options.Get.Path), options.Get.Name)
		case "show":
			return DisplaySharedFolder(options.Show.User, options.Show.Name)
		case "set":
			return UpdateSharedFolder(options.Set.User, options.Set.Name, options.Set.Set)
		default:
			return false
		}
	})
}
