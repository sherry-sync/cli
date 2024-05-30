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
		Name: helpers.If(yes, func() string {
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
		p = path.Join(".", helpers.If(helpers.IsUsernameFolder(name) == nil, func() string {
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

func responseToSource(response *api.ResponseFolder, userId string) config.Source {
	return config.Source{
		Id:          response.SherryId,
		Name:        response.Name,
		Access:      "write", // TODO: update with folder permissions
		OwnerId:     response.UserId,
		UserId:      userId,
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
}

func generateHashId(userId, sherryId string) string {
	return fmt.Sprintf("%s_%s_%d", userId, sherryId, time.Now().Unix())
}

func createWatcher(sourceId, userId, sherryId string, path string) config.Watcher {
	return config.Watcher{
		Source:    sourceId,
		LocalPath: path,
		HashesId:  generateHashId(userId, sherryId),
		UserId:    userId,
		Complete:  false,
	}
}

func generateSourceId(userId, sherryId string) string {
	return fmt.Sprintf("%s@%s", userId, sherryId)
}

func CreateSharedFolder(user string, yes bool, path string, name string, settings map[string]string) bool {
	credentials := auth.FindUserByUsername(user, true)

	if credentials == nil {
		helpers.PrintErr("User not found")
		return false
	}

	folderInfo := getFolderInfo(yes, path, name, settings)

	path = helpers.PreparePath(folderInfo.Path)

	for _, w := range config.GetConfig().Watchers {
		isChild, err := helpers.IsChildPath(path, helpers.PreparePath(w.LocalPath))
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
	sourceId := generateSourceId(credentials.UserId, response.SherryId)
	conf.Sources[sourceId] = responseToSource(response, credentials.UserId)
	conf.Watchers = append(conf.Watchers, createWatcher(sourceId, credentials.UserId, response.SherryId, path))

	helpers.PrintMessage(fmt.Sprintf("Sherry is created and watching at %s", path))

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

		source := helpers.Find(*availableFolders, func(f ResponseFolder) bool {
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

	// TODO: Create watcher with correct permissions

	return true
}

func DisplaySharedFolder(user string, name string) bool {
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

	for _, s := range *availableFolders {
		if s.Name != name {
			continue
		}
		source := responseToSource(&s, credentials.UserId)

		helpers.PrintMessage(fmt.Sprintf("Folder: %s", source.Name))
		helpers.PrintMessage("\n")
		helpers.PrintJson(source)
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

	source := helpers.Find(*availableFolders, func(f ResponseFolder) bool {
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
	estSource := responseToSource(response, credentials.UserId)
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

	helpers.PrintMessage(fmt.Sprintf("Folder was updated:"))
	helpers.PrintJson(estSource)

	return true
}

func UnwatchSharedFolder(path string, yes bool, force bool) bool {
	path = helpers.PreparePath(path)

	var watcher *config.Watcher
	for _, w := range config.GetConfig().Watchers {
		wPath := helpers.PreparePath(w.LocalPath)
		isChild, err := helpers.IsChildPath(path, wPath)
		if err != nil {
			helpers.PrintErr("Error while checking path")
			return false
		}
		if isChild && wPath != path {
			f := false
			if yes || helpers.Confirmation("Looks like it is not th root of shared directory, unwatch anyway?", "", &f) {
				watcher = &w
			} else {
				helpers.PrintErr("Aborting...")
				return false
			}
		} else {
			watcher = &w
		}
	}

	if watcher == nil {
		helpers.PrintErr("No watcher found")
		return false
	}

	conf := config.GetConfig()

	var newWatchers []config.Watcher
	for _, w := range conf.Watchers {
		if w.LocalPath != watcher.LocalPath {
			newWatchers = append(newWatchers, w)
		}
	}
	conf.Watchers = newWatchers

	// TODO: Add force option

	return true
}

func ListSharedFolders(user string, available bool) bool {
	var users []config.Credentials
	if user == "" {
		for _, c := range config.GetAuthConfig().Sources {
			users = append(users, c)
		}
	} else {
		credentials := auth.FindUserByUsername(user, false)
		if credentials == nil {
			helpers.PrintErr("User not found")
			return false
		}
		users = append(users, *credentials)
	}

	for _, u := range users {
		type Map struct {
			watchers []config.Watcher
			source   config.Source
		}
		var sources []config.Source
		if available {
			availableFolders, err := api.FolderGetAvailable(u.AccessToken)
			if err != nil {
				return false
			}
			for _, s := range *availableFolders {
				sources = append(sources, responseToSource(&s, u.UserId))
			}
		} else {
			for _, s := range config.GetConfig().Sources {
				if s.UserId == u.UserId {
					sources = append(sources, s)
				}
			}
		}
		helpers.PrintMessage(fmt.Sprintf("Folders for user: %s", auth.GetUserString(u)))
		helpers.PrintMessage("")
		for _, s := range sources {
			helpers.PrintMessage(fmt.Sprintf("Source: %s", s.Name))
			helpers.PrintJson(s)
			sourceKey := generateSourceId(u.UserId, s.Id)
			watchers := helpers.Filter(config.GetConfig().Watchers, func(watcher config.Watcher) bool {
				return sourceKey == watcher.Source
			})
			if len(watchers) != 0 {
				helpers.PrintMessage("Watching in these paths:")
				for _, w := range watchers {
					helpers.PrintMessage(w.LocalPath)
				}
			} else {
				helpers.PrintMessage("No currently watching paths")
			}
			helpers.PrintMessage("")
		}
		helpers.PrintMessage("")
	}

	return false
}
