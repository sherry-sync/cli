package folder

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/iancoleman/strcase"
	"os"
	"path"
	"sherry/shr/api"
	"sherry/shr/auth"
	"sherry/shr/config"
	"sherry/shr/constants"
	"sherry/shr/helpers"
	"strings"
	"time"
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

type ResponseFolderAllowedFileNames = struct {
	FileNameId string `json:"fileNameId"`
	Name       string `json:"name"`
	SherryId   string `json:"sherryId"`
}

type ResponseFolderAllowedFileTypes = struct {
	FileTypeId string `json:"fileTypeId"`
	Type       string `json:"type"`
	SherryId   string `json:"sherryId"`
}

type ResponseFolder = struct {
	SherryId         string                           `json:"sherryId"`
	Name             string                           `json:"name"`
	MaxFileSize      uint64                           `json:"maxFileSize"`
	MaxDirSize       uint64                           `json:"maxDirSize"`
	UserId           string                           `json:"userId"`
	AllowDir         bool                             `json:"allowDir"`
	AllowedFileTypes []ResponseFolderAllowedFileTypes `json:"allowedFileTypes"`
	AllowedFileNames []ResponseFolderAllowedFileNames `json:"allowedFileNames"`
}

func prepareSettings(settings map[string]string) map[string]string {
	s := make(map[string]string)

	for k, v := range settings {
		s[strcase.ToLowerCamel(k)] = v
	}

	return s
}

func getFolderSettings(yes bool, settings map[string]string) SourceSettings {
	settings = prepareSettings(settings)

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

	path = helpers.NormalizePath(folderInfo.Path)

	for _, w := range config.GetConfig().Watchers {
		isChild, err := helpers.IsChildPath(path, helpers.NormalizePath(w.LocalPath))
		if err != nil {
			helpers.PrintErr("Error while checking path")
			return false
		}
		if isChild {
			helpers.PrintErr("Path is already being watched")
			return false
		}
	}

	stat, err := os.Stat(path)
	if stat != nil && !stat.IsDir() {
		helpers.PrintErr("Path is not a directory")
		return false
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			helpers.PrintErr("Can't create directory")
			helpers.PrintErr(err.Error())
			return false
		}
	}

	response, err := api.FolderCreate(Payload{
		Name:             folderInfo.Name,
		AllowDir:         folderInfo.Settings.AllowDir,
		MaxFileSize:      folderInfo.Settings.MaxFileSize,
		MaxDirSize:       folderInfo.Settings.MaxDirSize,
		AllowedFileNames: folderInfo.Settings.AllowedFileNames,
		AllowedFileTypes: folderInfo.Settings.AllowedFileTypes,
	}, credentials.AccessToken)
	if err != nil {
		return false
	}

	conf := config.GetConfig()
	sourceId := fmt.Sprintf("%s@%s", credentials.UserId, response.SherryId)
	conf.Sources[sourceId] = config.Source{
		Id:          response.SherryId,
		Name:        response.Name,
		Access:      "write",
		OwnerId:     response.UserId,
		UserId:      credentials.UserId,
		AllowDir:    response.AllowDir,
		MaxFileSize: response.MaxFileSize,
		MaxDirSize:  response.MaxDirSize,
		AllowedFileNames: helpers.Map(response.AllowedFileNames, func(f ResponseFolderAllowedFileNames) string {
			return f.Name
		}),
		AllowedFileTypes: helpers.Map(response.AllowedFileTypes, func(f ResponseFolderAllowedFileTypes) string {
			return f.Type
		}),
	}
	conf.Watchers = append(conf.Watchers, config.Watcher{
		Source:    sourceId,
		LocalPath: path,
		HashesId:  fmt.Sprintf("%s_%s_%d", credentials.UserId, response.SherryId, time.Now().Unix()),
		UserId:    credentials.UserId,
		Complete:  false,
	})

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

	var folderId string
	if helpers.IsUsernameFolder(folderParams.Name) == nil {
		availableFolders, err := api.FolderGetAvailable(credentials.AccessToken)
		if err != nil {
			return false
		}

		args := strings.Split(folderParams.Name, ":")
		folderName := args[1]
		userData, err := api.UserFind(args[0], credentials.AccessToken)
		if err != nil {
			return false
		}

		source := helpers.FindIn(*availableFolders, func(f ResponseFolder) bool {
			return f.Name == folderName && f.UserId == userData.UserId
		})
		if source == nil {
			helpers.PrintErr("Folder is not available or not exists")
			return false
		}
		folderId = source.SherryId
	} else {
		folderId = folderParams.Name
	}

	folder, err := api.FolderGet(folderId, credentials.AccessToken)
	if err != nil {
		return false
	}

	helpers.PrintJson(folder)

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

	availableFolders, err := api.FolderGetAvailable(credentials.AccessToken)
	if err != nil {
		return false
	}

	source := helpers.FindIn(*availableFolders, func(f ResponseFolder) bool {
		return f.Name == name && f.UserId == credentials.UserId
	})
	if source == nil {
		helpers.PrintErr("Folder is not available or not exists")
		return false
	}

	response, err := api.FolderUpdate(source.SherryId, Payload{
		AllowDir:    helpers.ParseBool("Allow directory", settings["allowDir"], source.AllowDir),
		MaxFileSize: helpers.ParseDataSize("Max file size", settings["maxFileSize"], source.MaxFileSize, constants.MaxFileSize),
		MaxDirSize:  helpers.ParseDataSize("Max directory size", settings["maxDirSize"], source.MaxDirSize, constants.MaxDirSize),
		AllowedFileNames: helpers.ParseValueArray(
			"Allowed file names",
			settings["allowedFileNames"],
			helpers.IsGlobValidator, helpers.ToJoinedValues(helpers.Map(source.AllowedFileNames, func(f ResponseFolderAllowedFileNames) string {
				return f.Name
			})),
		),
		AllowedFileTypes: helpers.ParseValueArray(
			"Allowed file types",
			settings["allowedFileTypes"],
			helpers.IsMimeTypeValidator, helpers.ToJoinedValues(helpers.Map(source.AllowedFileTypes, func(f ResponseFolderAllowedFileTypes) string {
				return f.Type
			})),
		),
	}, credentials.AccessToken)
	if err != nil {
		return false
	}

	conf := config.GetConfig()
	estSource := config.Source{
		AllowDir:    response.AllowDir,
		MaxFileSize: response.MaxFileSize,
		MaxDirSize:  response.MaxDirSize,
		AllowedFileNames: helpers.Map(response.AllowedFileNames, func(f ResponseFolderAllowedFileNames) string {
			return f.Name
		}),
		AllowedFileTypes: helpers.Map(response.AllowedFileTypes, func(f ResponseFolderAllowedFileTypes) string {
			return f.Type
		}),
	}
	for _, s := range conf.Sources {
		if s.Id != source.SherryId {
			continue
		}

		s.AllowDir = estSource.AllowDir
		s.MaxFileSize = estSource.MaxFileSize
		s.MaxDirSize = estSource.MaxDirSize
		s.AllowedFileNames = estSource.AllowedFileNames
		s.AllowedFileTypes = estSource.AllowedFileTypes
	}

	return true
}
