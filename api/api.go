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
	"strings"
)

type ErrorResponse = struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type ErrorResponseArray = struct {
	Message    []string `json:"message"`
	StatusCode int      `json:"statusCode"`
}

func getUrl(route string) string {
	base, err := url.Parse(config.GetConfig().ApiUrl)
	if err != nil {
		helpers.PrintErr("Can't parse API URL")
		return ""
	}
	parts := strings.SplitN(route, "?", 2)
	base.Path = path.Join(base.Path, parts[0])
	if len(parts) == 2 {
		base.RawQuery = parts[1]
	}
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
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, getUrl(route), nil)
	} else {
		req, err = http.NewRequest(method, getUrl(route), body)
	}
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

func Delete(route string, auth string) (string, error) {
	return authRequest(http.MethodDelete, route, nil, auth)
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
			if json.Unmarshal([]byte(res), &resErr) == nil {
				helpers.PrintErr(resErr.Message)
			} else {
				var resErrArr ErrorResponseArray
				err := json.Unmarshal([]byte(res), &resErrArr)
				if err != nil {
					helpers.PrintErr(res)
				} else {
					helpers.PrintErr("Couple errors found:")
					for _, m := range resErrArr.Message {
						helpers.PrintErr(fmt.Sprintf("  %s", m))
					}
				}
			}

		} else {
			helpers.PrintErr(err.Error())
		}
		return "", err
	}
	return res, nil
}

func ParseResponse[T any](res string) (*T, error) {
	var v T
	if err := json.Unmarshal([]byte(res), &v); err != nil {
		helpers.PrintErr("Can't parse response")
		helpers.PrintErr(err.Error())
		return nil, err
	}
	return &v, nil
}
