package config

import (
	"encoding/json"
	"errors"
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
	ApiUrl   string            `json:"apiUrl"`
	Sources  map[string]Source `json:"sources"`
	Watchers []Watcher         `json:"watchers"`
	Webhooks []string          `json:"webhooks"`
}

type Credentials struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthorizationConfig struct {
	Sources map[string]Credentials `json:"records"`
	Default string                 `json:"default"`
}

var configPath = ""
var globalConfig *Config = nil
var globalAuthConfig *AuthorizationConfig = nil

func GetConfigPath(overwritePath string) string {
	configPath := overwritePath
	if configPath == "" {
		home, _ := os.UserHomeDir()
		configPath = fmt.Sprintf(path.Join(home, constants.ConfigDir))
	}
	return configPath
}

func SetConfig(c *Config) {
	globalConfig = c
}

func SetAuthConfig(c *AuthorizationConfig) {
	globalAuthConfig = c
}

func ReadConfig() *Config {
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

func ReadAuthConfig() *AuthorizationConfig {
	file, err := os.ReadFile(path.Join(configPath, constants.AuthConfigFile))

	if err != nil {
		helpers.PrintErr(fmt.Sprintf("Can't find authorization configuration, searching \"%s\"", configPath))
		return nil
	}

	var c AuthorizationConfig
	if err := json.Unmarshal(file, &c); err != nil {
		helpers.PrintErr(fmt.Sprintf("Unable to parse authorization configuration file: %s", err))
		return nil
	}

	return &c
}

func SetupConfig(overwritePath string) error {
	configPath = GetConfigPath(overwritePath)

	c := ReadConfig()
	if c == nil {
		return errors.New("can't find configuration")
	}
	SetConfig(c)

	auth := ReadAuthConfig()
	if auth == nil {
		return errors.New("can't find authorization configuration")
	}
	SetAuthConfig(auth)

	return nil
}

func GetConfig() *Config {
	return globalConfig
}

func GetAuthConfig() *AuthorizationConfig {
	return globalAuthConfig
}

func CommitConfig() {
	data, _ := json.MarshalIndent(globalConfig, "", "  ")
	err := os.WriteFile(path.Join(configPath, constants.ConfigFile), data, 0644)
	if err != nil {
		helpers.PrintErr(fmt.Sprintf("Unable to save configuration: %s", err))
		return
	}
}

func CommitAuth() {
	data, _ := json.MarshalIndent(globalAuthConfig, "", "  ")
	err := os.WriteFile(path.Join(configPath, constants.AuthConfigFile), data, 0644)
	if err != nil {
		helpers.PrintErr(fmt.Sprintf("Unable to save authorization configuration: %s", err))
		return
	}
}
