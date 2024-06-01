package folder

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
)

type Options struct {
	Create      CreateOptions     `command:"create" description:"Create shared folder"`
	Get         GetOptions        `command:"get" description:"Get shared folder"`
	Show        ShowOptions       `command:"show" description:"Display shared folder sharing info"`
	Update      UpdateOptions     `command:"update" description:"Set shared folder settings"`
	Permissions PermissionOptions `command:"permission" description:"Manage shared folder access"`
	List        ListOptions       `command:"list" description:"List folders"`
	Unwatch     UnwatchOptions    `command:"unwatch" description:"Unwatch folder"`
}

type UnwatchOptions struct {
	Yes   bool `long:"yes" short:"y" description:"Skip confirmation and use default values where possible"`
	Force bool `long:"force" short:"f" description:"Also disable folder shearing"`
	Args  struct {
		Path flag.Filename `positional-arg-name:"path"  description:"Shared folder path"`
	} `positional-args:"yes" description:"Shared folder name"`
}

type ListOptions struct {
	User      string `long:"user" short:"u" description:"Choose specific user"`
	Available bool   `long:"available" short:"a" description:"Show all available folders"`
}

type PermissionGrantOptions struct {
	User   string `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Target string `long:"target" short:"t" description:"Username or id of user to manage access for"`
	Role   string `long:"role" description:"Permission role (read/write)"`
	Name   string `long:"name" short:"n" description:"Shared folder name"`
}

type PermissionRevokeOptions struct {
	User   string `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Target string `long:"target" short:"t" description:"Username or id of user to manage access for"`
	Name   string `long:"name" short:"n" description:"Shared folder name"`
}

type PermissionOptions struct {
	Grant  PermissionGrantOptions  `command:"grant" description:"Grant access to shared folder"`
	Revoke PermissionRevokeOptions `command:"revoke" description:"Revoke access from shared folder"`
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
		case "unwatch":
			return UnwatchSharedFolder(string(options.Unwatch.Args.Path), options.Unwatch.Yes, options.Unwatch.Force)
		case "list":
			return ListSharedFolders(options.List.User, options.List.Available)
		case "permission":
			switch cmd.Active.Active.Active.Name {
			case "grant":
				return GrantPermission(options.Permissions.Grant.User, options.Permissions.Grant.Target, options.Permissions.Grant.Name, options.Permissions.Grant.Role)
			case "revoke":
				return RevokePermission(options.Permissions.Revoke.User, options.Permissions.Revoke.Target, options.Permissions.Revoke.Name)
			default:
				return false
			}
		default:
			return false
		}
	})
}
