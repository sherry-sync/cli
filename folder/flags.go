package folder

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
)

type Options struct {
	Create CreateOptions `command:"create" description:"Create shared folder"`
	Get    GetOptions    `command:"get" description:"Get shared folder"`
	Show   ShowOptions   `command:"show" description:"Display shared folder sharing info"`
	Update UpdateOptions `command:"update" description:"Set shared folder settings"`
	Access AccessOptions `command:"access" description:"Manage shared folder access"`
}

type AccessOptions struct {
	User   string `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Name   string `long:"name" short:"n" description:"Specify shared folder name"`
	Action string `long:"action" short:"a" description:"Specify action to perform (GRANT/REFUSE)"`
	Args   struct {
		Folder string `positional-arg-name:"user"  description:"Username or id of user to manage access for"`
	} `positional-args:"yes" description:"Username or id of user to manage access for"`
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
	Yes  bool          `long:"yes" short:"y" description:"Skip confirmation and use default values where possible"`
	Args struct {
		Folder string `positional-arg-name:"folder"  description:"Shared folder in format owner_username:folder_name or folder id"`
	} `positional-args:"yes" required:"yes" description:"Shared folder in format owner_username:folder_name or folder id"`
}

type ShowOptions struct {
	User string `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Args struct {
		Name string `positional-arg-name:"folder"  description:"Shared folder name"`
	} `positional-args:"yes" required:"yes" description:"Shared folder name"`
}

type UpdateOptions struct {
	User string            `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Set  map[string]string `long:"set" short:"s" description:"Set folder settings"`
	Args struct {
		Name string `positional-arg-name:"folder"  description:"Shared folder name"`
	} `positional-args:"yes" required:"yes" description:"Shared folder name"`
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
			return GetSharedFolder(options.Get.User, options.Get.Yes, string(options.Get.Path), options.Get.Args.Folder)
		case "show":
			return DisplaySharedFolder(options.Show.User, options.Show.Args.Name)
		case "update":
			return UpdateSharedFolder(options.Update.User, options.Update.Args.Name, options.Update.Set)
		default:
			return false
		}
	})
}
