package main

import (
	flag "github.com/jessevdk/go-flags"
	"sherry/shr/config"
)

var parser = flag.NewParser(&defaultFlags, flag.Default)

func init() {
	initCommand(parser.Command)
}

func main() {
	if _, err := parser.Parse(); err != nil {
		return
	}

	c := config.ReadConfig(string(defaultFlags.ConfigPath))
	if c == nil {
		return
	}
	applyCommand(parser.Command, *c)
}
