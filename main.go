package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/jessevdk/go-flags"
	"os"
	"path"
)

func printErr(err string) {
	_, _ = fmt.Fprintf(os.Stderr, err)
}

var parser = flag.NewParser(&defaultFlags, flag.Default)

func init() {
	initFlags(parser.Command)
}

func main() {
	if _, err := parser.Parse(); err != nil {
		return
	}

	applyCommand(parser.Command)
	printJson(defaultFlags)
	printJson(authFlags)
	printJson(folderFlags)

	configPath := string(defaultFlags.ConfigPath)
	if configPath == "" {
		home, _ := os.UserHomeDir()
		configPath = fmt.Sprintf(path.Join(home, ConfigDir))
	}

	file, err := os.ReadFile(path.Join(configPath, ConfigFile))

	if err != nil {
		printErr(fmt.Sprintf("Can't find configuration, searching \"%s\"", configPath))
		return
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		printErr(fmt.Sprintf("Unable to parse configuration file: %s", err))
		return
	}

	printJson(config)
}
