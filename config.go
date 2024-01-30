package main

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
