package folder

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"sherry/shr/api"
	"sherry/shr/auth"
	"sherry/shr/constants"
	"sherry/shr/helpers"
)

type SourceSettings struct {
	AllowDir         bool     `json:"allowDir"`
	MaxFileSize      uint64   `json:"maxFileSize"`
	MaxDirSize       uint64   `json:"maxDirSize"`
	AllowedFileNames []string `json:"allowedFileNames"`
	AllowedFileTypes []string `json:"allowedFileTypes"`
}

type Info = struct {
	Name     string         `json:"name"`
	Path     string         `json:"path"`
	Settings SourceSettings `json:"settings"`
}

type Payload = struct {
	Name             string   `json:"name"`
	AllowDir         bool     `json:"allowDir"`
	MaxFileSize      uint64   `json:"maxFileSize"`
	MaxDirSize       uint64   `json:"maxDirSize"`
	AllowedFileNames []string `json:"allowedFileNames"`
	AllowedFileTypes []string `json:"allowedFileTypes"`
}

func getFolderSettings(yes bool, settings map[string]string) SourceSettings {
	if yes {
		return SourceSettings{
			AllowDir:         helpers.ParseBool("Allow directory", settings["allowDir"], constants.DefaultAllowDir),
			MaxFileSize:      helpers.ParseDataSize("Max file size", settings["maxFileSize"], constants.DefaultMaxFileSize, constants.MaxFileSize),
			MaxDirSize:       helpers.ParseDataSize("Max directory size", settings["maxDirSize"], constants.DefaultMaxDirSize, constants.MaxDirSize),
			AllowedFileNames: helpers.ParseValueArray("Allowed file names", settings["allowedFileNames"], helpers.IsGlobValidator),
			AllowedFileTypes: helpers.ParseValueArray("Allowed file types", settings["allowedFileTypes"], helpers.IsMimeTypeValidator),
		}
	}

	defaultAllowDir := new(bool)
	*defaultAllowDir = constants.DefaultAllowDir

	return SourceSettings{
		AllowDir: helpers.Confirmation("Allow directory", settings["allowDir"], defaultAllowDir),
		MaxFileSize: helpers.ParseDataSize(
			"Max file size",
			helpers.Input(
				"Max file size",
				settings["maxFileSize"],
				helpers.GetDataSizeValidator(constants.MaxFileSize),
				fmt.Sprintf("Up to %s", humanize.Bytes(constants.MaxFileSize)),
				false,
			),
			constants.DefaultMaxFileSize,
			constants.MaxFileSize,
		),
		MaxDirSize: helpers.ParseDataSize(
			"Max directory size",
			helpers.Input(
				"Max directory size",
				settings["maxDirSize"],
				helpers.GetDataSizeValidator(constants.MaxDirSize),
				fmt.Sprintf("Up to %s", humanize.Bytes(constants.MaxDirSize)),
				false,
			),
			constants.DefaultMaxDirSize,
			constants.MaxDirSize,
		),
		AllowedFileNames: helpers.ParseValueArray(
			"Allowed file names",
			helpers.Input(
				"Allowed file names",
				settings["allowedFileNames"],
				helpers.GetValidValueArrayValidator(true, helpers.IsGlobValidator),
				"Glob patterns separated by commas",
				false,
			),
			helpers.IsGlobValidator,
		),
		AllowedFileTypes: helpers.ParseValueArray(
			"Allowed file types",
			helpers.Input(
				"Allowed file types",
				settings["allowedFileTypes"],
				helpers.GetValidValueArrayValidator(true, helpers.IsMimeTypeValidator),
				"MIME types separated by commas",
				false,
			),
			helpers.IsMimeTypeValidator,
		),
	}
}

func getFolderInfo(yes bool, path string, settings map[string]string) Info {
	if yes {
		return Info{
			Name:     "",
			Path:     helpers.Input("Path", path, helpers.IsPathValidator, "", false),
			Settings: getFolderSettings(yes, settings),
		}
	}

	return Info{
		Name:     helpers.Input("Name", "", helpers.IsWordValidator, "", false),
		Path:     helpers.Input("Path", path, helpers.IsPathValidator, "", false),
		Settings: getFolderSettings(yes, settings),
	}
}

func CreateSharedFolder(user string, yes bool, path string, name string, settings map[string]string) bool {
	var credentials = auth.FindUserByUsername(user, true)

	if credentials == nil {
		helpers.PrintErr("User not found")
		return false
	}

	folderInfo := getFolderInfo(yes, path, settings)

	// TODO: Validate path

	body, _ := json.Marshal(Payload{
		Name:             folderInfo.Name,
		AllowDir:         folderInfo.Settings.AllowDir,
		MaxFileSize:      folderInfo.Settings.MaxFileSize,
		MaxDirSize:       folderInfo.Settings.MaxDirSize,
		AllowedFileNames: folderInfo.Settings.AllowedFileNames,
		AllowedFileTypes: folderInfo.Settings.AllowedFileTypes,
	})
	res, err := api.Post("/sherry", body, credentials.AccessToken)
	if _, err := helpers.ValidateResponse(res, err); err != nil {
		return false
	}

	return true
}
