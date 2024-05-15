package main

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
)

func main() {
	var options Options
	var parser = flag.NewParser(&options, flag.Default)

	if _, err := parser.Parse(); err != nil {
		return
	}

	c := config.SetupConfig(string(options.ConfigPath))
	if c != nil {
		return
	}
	applyCommand(parser.Command, options)
}
