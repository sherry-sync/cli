package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sherry/shr/constants"
	"sherry/shr/helpers"
)

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
	HashesId  string `json:"hashesId"`
	UserId    string `json:"userId"`
}

type Config struct {
	Sources  map[string]Source `json:"sources"`
	Watchers []Watcher         `json:"watchers"`
	Webhooks []string          `json:"webhooks"`
}

func GetConfigPath(overwritePath string) string {
	configPath := overwritePath
	if configPath == "" {
		home, _ := os.UserHomeDir()
		configPath = fmt.Sprintf(path.Join(home, constants.ConfigDir))
	}
	return configPath
}

func ReadConfig(overwritePath string) *Config {
	configPath := GetConfigPath(overwritePath)
	file, err := os.ReadFile(path.Join(configPath, constants.ConfigFile))

	if err != nil {
		helpers.PrintErr(fmt.Sprintf("Can't find configuration, searching \"%s\"", configPath))
		return nil
	}

	var c Config
	if err := json.Unmarshal(file, &c); err != nil {
		helpers.PrintErr(fmt.Sprintf("Unable to parse configuration file: %s", err))
		return nil
	}

	return &c
}
