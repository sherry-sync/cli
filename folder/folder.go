package folder

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"path"
	"sherry/shr/api"
	"sherry/shr/auth"
	"sherry/shr/config"
	"sherry/shr/constants"
	"sherry/shr/helpers"
	"strings"
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

type Params = struct {
	Name string `json:"name"`
	Path string `json:"path"`
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
			AllowedFileNames: helpers.ParseValueArray("Allowed file names", settings["allowedFileNames"], helpers.IsGlobValidator, ""),
			AllowedFileTypes: helpers.ParseValueArray("Allowed file types", settings["allowedFileTypes"], helpers.IsMimeTypeValidator, ""),
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
			"",
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
			"",
		),
	}
}

func getFolderInfo(yes bool, path string, name string, settings map[string]string) Info {
	return Info{
		Name: helpers.IfLazy(yes, func() string {
			return name
		}, func() string {
			return helpers.Input("Name", name, helpers.IsWordValidator, "", false)
		}),
		Path:     helpers.Input("Path", path, helpers.IsPathValidator, "", false),
		Settings: getFolderSettings(yes, settings),
	}
}

func getFolderParams(yes bool, p string, name string) Params {
	name = helpers.Input("Folder name in format owner_username:folder_name or id", name, helpers.IsUsernameFolderOrId, "", false)

	if yes && p == "" {
		p = path.Join(".", helpers.IfLazy(helpers.IsUsernameFolder(name) == nil, func() string {
			return strings.Split(name, ":")[1]
		}, func() string {
			return name
		}))
	}

	return Params{
		Name: name,
		Path: helpers.Input("Path", p, helpers.IsPathValidator, "", false),
	}
}

func CreateSharedFolder(user string, yes bool, path string, name string, settings map[string]string) bool {
	credentials := auth.FindUserByUsername(user, true)

	if credentials == nil {
		helpers.PrintErr("User not found")
		return false
	}

	folderInfo := getFolderInfo(yes, path, name, settings)

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
	if _, err := api.ValidateResponse(res, err); err != nil {
		return false
	}

	return true
}

func GetSharedFolder(user string, yes bool, path string, name string) bool {
	credentials := auth.FindUserByUsername(user, true)

	if credentials == nil {
		helpers.PrintErr("User not found")
		return false
	}

	folderParams := getFolderParams(yes, path, name)

	if helpers.IsExists(folderParams.Path) {
		helpers.PrintErr("Directory already exists")
		return false
	}

	res, err := api.Get(fmt.Sprintf("/sherry/%s", folderParams.Name), credentials.AccessToken)
	if _, err := api.ValidateResponse(res, err); err != nil {
		return false
	}

	// TODO: Download files

	return true
}

func DisplaySharedFolder(user string, name string) bool {
	name = helpers.Input("Folder name", name, helpers.IsWordValidator, "", false)

	credentials := auth.FindUserByUsername(user, false)
	if credentials == nil && user != "" {
		helpers.PrintErr("User not found")
		return false
	}

	sources := config.GetConfig().Sources

	for _, s := range sources {
		if user != "" && s.Name != name {
			continue
		}
		helpers.PrintMessage(fmt.Sprintf("Folder: %s", s.Name))
		helpers.PrintMessage(fmt.Sprintf("User: %s", helpers.If(credentials != nil, func(c bool) string {
			if c {
				return credentials.Username
			} else {
				return auth.GetUserById(s.UserId).Username
			}
		})))
		helpers.PrintMessage("\n")
		helpers.PrintJson(s)
	}

	return false
}

func UpdateSharedFolder(user string, name string, settings map[string]string) bool {
	name = helpers.Input("Folder name", name, helpers.IsWordValidator, "", false)
	credentials := auth.FindUserByUsername(user, true)
	if credentials == nil {
		helpers.PrintErr("User not found")
		return false
	}

	source, _ := helpers.Find(config.GetConfig().Sources, func(source config.Source) bool {
		return source.Name == name && source.OwnerId == credentials.UserId
	})
	if source == nil {
		helpers.PrintErr("Source not found")
		return false
	}

	body, _ := json.Marshal(Payload{
		AllowDir:         helpers.ParseBool("Allow directory", settings["allowDir"], source.AllowDir),
		MaxFileSize:      helpers.ParseDataSize("Max file size", settings["maxFileSize"], source.MaxFileSize, constants.MaxFileSize),
		MaxDirSize:       helpers.ParseDataSize("Max directory size", settings["maxDirSize"], source.MaxDirSize, constants.MaxDirSize),
		AllowedFileNames: helpers.ParseValueArray("Allowed file names", settings["allowedFileNames"], helpers.IsGlobValidator, helpers.ToJoinedValues(source.AllowedFileNames)),
		AllowedFileTypes: helpers.ParseValueArray("Allowed file types", settings["allowedFileTypes"], helpers.IsMimeTypeValidator, helpers.ToJoinedValues(source.AllowedFileTypes)),
	})
	res, err := api.Patch(fmt.Sprintf("/sherry/%s", source.Id), body, credentials.AccessToken)
	if _, err := api.ValidateResponse(res, err); err != nil {
		return false
	}

	// TODO: update config

	return true
}
