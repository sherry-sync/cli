package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type PayloadFolder = struct {
	Name             string   `json:"name"`
	AllowDir         bool     `json:"allowDir"`
	MaxFileSize      uint64   `json:"maxFileSize"`
	MaxDirSize       uint64   `json:"maxDirSize"`
	AllowedFileNames []string `json:"allowedFileNames"`
	AllowedFileTypes []string `json:"allowedFileTypes"`
}

type PayloadFolderPermissionRole = string

const (
	PermissionRoleOwner PayloadFolderPermissionRole = "OWNER"
	PermissionRoleRead  PayloadFolderPermissionRole = "READ"
	PermissionRoleWrite PayloadFolderPermissionRole = "WRITE"
)

type PayloadFolderPermissionAction = string

const (
	PermissionActionGrant  PayloadFolderPermissionAction = "GRANT"
	PermissionActionRefuse PayloadFolderPermissionAction = "REFUSE"
)

type PayloadFolderPermission = struct {
	Role   PayloadFolderPermissionRole   `json:"role"`
	Action PayloadFolderPermissionAction `json:"action"`
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

type SherryPermission = struct {
	SherryPermissionId string `json:"sherryPermissionId"`
	Role               string `json:"role"`
	SherryId           string `json:"sherryId"`
	UserId             string `json:"userId"`
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
	SherryPermission []SherryPermission               `json:"sherryPermission"`
}

type FileType = string

const (
	Dir  FileType = "DIR"
	File FileType = "FILE"
)

type FileResponse struct {
	SherryFileID string   `json:"sherryFileId"`
	SherryID     string   `json:"sherryId"`
	Path         string   `json:"path"`
	OldPath      string   `json:"oldPath"`
	Hash         string   `json:"hash"`
	Size         uint64   `json:"size"`
	CreatedAt    uint64   `json:"createdAt"`
	UpdatedAt    uint64   `json:"updatedAt"`
	FileType     FileType `json:"fileType"`
}

func FolderCreate(payload PayloadFolder, accessToken string) (*ResponseFolder, error) {
	body, _ := json.Marshal(payload)
	res, err := ValidateResponse(Post("/sherry", body, accessToken))
	if err != nil {
		return nil, err
	}

	return ParseResponse[ResponseFolder](res)
}

func FolderUpdate(id string, payload PayloadFolder, accessToken string) (*ResponseFolder, error) {
	body, _ := json.Marshal(payload)
	res, err := ValidateResponse(Patch(fmt.Sprintf("/sherry/%s", id), body, accessToken))
	if err != nil {
		return nil, err
	}

	return ParseResponse[ResponseFolder](res)
}

func FolderGetAvailable(accessToken string) (*[]ResponseFolder, error) {
	res, err := ValidateResponse(Get("/sherry/my", accessToken))
	if err != nil {
		return nil, err
	}

	return ParseResponse[[]ResponseFolder](res)
}

func FolderGet(id string, accessToken string) (*ResponseFolder, error) {
	res, err := ValidateResponse(Get(fmt.Sprintf("/sherry/%s", id), accessToken))
	if err != nil {
		return nil, err
	}

	return ParseResponse[ResponseFolder](res)
}

func FolderDelete(id string, accessToken string) error {
	_, err := ValidateResponse(Delete(fmt.Sprintf("/sherry/%s", id), accessToken))
	if err != nil {
		return err
	}

	return nil
}

func FolderPermission(folderId, userId string, payload PayloadFolderPermission, accessToken string) error {
	body, _ := json.Marshal(payload)
	_, err := ValidateResponse(Patch(fmt.Sprintf("/sherry/%s/users/%s/permission", folderId, userId), body, accessToken))
	if err != nil {
		return err
	}

	return nil
}

func FolderFiles(id string, accessToken string) (*[]FileResponse, error) {
	res, err := ValidateResponse(Get(fmt.Sprintf("/file/%s", id), accessToken))
	if err != nil {
		return nil, err
	}

	return ParseResponse[[]FileResponse](res)
}

func FolderFileDownload(id, filePath string, accessToken string, dst string) error {
	req, err := http.NewRequest(http.MethodGet, getUrl(fmt.Sprintf("/file/instance/%s?path=%s", id, filePath)), nil)
	if err != nil {
		return err
	}

	// Set the Authorization header
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return err
}
