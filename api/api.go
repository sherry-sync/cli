package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sherry/shr/config"
	"sherry/shr/helpers"
)

type ErrorResponse = struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func getUrl(route string) string {
	base, err := url.Parse(config.GetConfig().ApiUrl)
	if err != nil {
		helpers.PrintErr("Can't parse API URL")
		return ""
	}
	base.Path = path.Join(base.Path, route)
	return base.String()
}

var UnsuccessfulResponseCodeError = errors.New("unsuccessful response code")

func isSuccess(res *http.Response) bool {
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return true
	}
	return false
}

func parse(res *http.Response) (string, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	str := string(body)
	if !isSuccess(res) {
		return str, UnsuccessfulResponseCodeError
	}
	return str, nil
}

func authRequest(method string, route string, body *bytes.Buffer, auth string) (string, error) {
	req, err := http.NewRequest(method, getUrl(route), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", fmt.Sprint("Bearer ", auth))
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	return parse(res)
}

func Get(route string, auth string) (string, error) {
	return authRequest(http.MethodGet, route, nil, auth)
}

func Post(route string, body []byte, auth string) (string, error) {
	return authRequest(http.MethodPost, route, bytes.NewBuffer(body), auth)
}

func Patch(route string, body []byte, auth string) (string, error) {
	return authRequest(http.MethodPatch, route, bytes.NewBuffer(body), auth)
}

func ValidateResponse(res string, err error) (string, error) {
	if err != nil {
		if res != "" {
			var resErr ErrorResponse
			err := json.Unmarshal([]byte(res), &resErr)
			if err != nil {
				helpers.PrintErr(res)
			}
			helpers.PrintErr(resErr.Message)
		} else {
			helpers.PrintErr(err.Error())
		}
		return "", err
	}
	return res, nil
}
