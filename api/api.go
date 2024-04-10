package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"sherry/shr/config"
	"sherry/shr/helpers"
)

func getUrl(route string) string {
	base, err := url.Parse(config.GetConfig().ApiUrl)
	if err != nil {
		helpers.PrintErr("Can't parse API URL")
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

func Get(route string) (string, error) {
	res, err := http.Get(getUrl(route))
	if err != nil {
		return "", err
	}
	return parse(res)
}

func Post(route string, body []byte) (string, error) {
	res, err := http.Post(getUrl(route), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	return parse(res)
}

func Patch(route string, body []byte) (string, error) {
	req, err := http.NewRequest(http.MethodPatch, getUrl(route), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	return parse(res)
}
