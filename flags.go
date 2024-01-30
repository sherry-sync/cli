package main

import (
	"fmt"
	flag "github.com/jessevdk/go-flags"
)

func AddCommand(cmd *flag.Command, command string, shortDescription string, longDescription string, data interface{}) {
	if _, err := cmd.AddCommand(
		command,
		shortDescription,
		longDescription,
		data,
	); err != nil {
		fmt.Println(err)
		return
	}
}

var defaultFlags struct {
	ConfigPath flag.Filename `long:"config" short:"c" description:"Path to configuration folder"`
}

var authCommandGroup = "auth"
var authFlags struct {
	Register bool   `long:"register" short:"r" description:"Register new user"`
	User     string `long:"user" short:"u" description:"Username"`
	Password string `long:"password" short:"p" description:"User password"`
	Default  bool   `long:"default" short:"d" description:"Get/Set default user (user should be specified to set)"`
}

func initAuthCommand(cmd *flag.Command) {
	AddCommand(cmd,
		authCommandGroup,
		"Authorization",
		"Manage your authorization",
		&authFlags,
	)
}

var folderCommandGroup = "folder"
var folderFlags struct {
	User   string            `long:"user" short:"u" description:"Use specific user profile for operation (Default profile will be used if no specified)"`
	Path   flag.Filename     `long:"path" short:"p" description:"Specify local path for operation"`
	Create bool              `long:"create" short:"c" description:"Create shared folder"`
	Get    bool              `long:"get" short:"g" description:"Get shared folder"`
	Set    map[string]string `long:"set" short:"s" description:"Set folder settings"`
}

func initFolderCommand(cmd *flag.Command) {
	AddCommand(cmd,
		folderCommandGroup,
		"Folder",
		"Manage your folders",
		&folderFlags,
	)
}

func initFlags(cmd *flag.Command) {
	initAuthCommand(cmd)
	initFolderCommand(cmd)
}

func applyCommand(_ *flag.Command) {

}
