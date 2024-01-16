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

const ConfigDir = ".sherry"
const ConfigFile = "config.json"

var opts struct {
	ConfigPath flag.Filename `long:"config" short:"c" description:"Path to configuration folder"`
}

type Source struct {
	Id               string   `json:"id"`
	Access           string   `json:"access"`
	MaxFileSize      uint64   `json:"maxFileSize"`
	MaxDirSize       uint64   `json:"maxDirSize"`
	AllowedFileNames []string `json:"allowedFileNames"`
	AllowedFileTypes []string `json:"allowedFileTypes"`
}

type Watcher struct {
	Source    string `json:"source"`
	LocalPath string `json:"localPath"`
}

type Config struct {
	Sources  map[string]Source `json:"sources"`
	Watchers []Watcher         `json:"watchers"`
	Webhooks []string          `json:"webhooks"`
}

func main() {
	if _, err := flag.Parse(&opts); err != nil {
		return
	}

	configPath := string(opts.ConfigPath)
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

	fmt.Println(config)
}
