package api

import (
	"encoding/json"
	"fmt"
)

type PayloadFolder = struct {
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
