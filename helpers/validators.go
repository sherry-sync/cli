package helpers

import (
	"github.com/dlclark/regexp2"
	"github.com/dustin/go-humanize"
	"github.com/erikgeiser/promptkit/textinput"
	"mime"
	"net/mail"
	"path/filepath"
)

var isPasswordRegex = regexp2.MustCompile(`(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}`, 0).MatchString
var isWordRegex = regexp2.MustCompile(`^\w+$`, 0).MatchString
var isIdRegex = regexp2.MustCompile(`^([a-zA-Z0-9]+-)+[a-zA-Z0-9]+$`, 0).MatchString
var isUsernameFolder = regexp2.MustCompile(`^\w+:\w+$`, 0).MatchString
var isWindowsPathBanChar = regexp2.MustCompile(`[/<>:"\|?*\x00\x08\x0B\x0C\x0E-\x1F]`, 0).MatchString
var isLinuxPathBanChar = regexp2.MustCompile(`[/\x00]`, 0).MatchString
var isWindowsPathBanName = regexp2.MustCompile(`^(CON|PRN|AUX|NUL|COM1|COM2|COM3|COM4|COM5|COM6|COM7|COM8|COM9|LPT1|LPT2|LPT3|LPT4|LPT5|LPT6|LPT7|LPT8|LPT9)(\..*)?$`, regexp2.IgnoreCase).MatchString
var isWindowsPathBanEnd = regexp2.MustCompile(`[. ]$`, 0).MatchString
var isWindowsDrive = regexp2.MustCompile(`^[a-zA-Z]:$`, 0).MatchString

func match(regex func(s string) (bool, error), input string) bool {
	match, _ := regex(input)
	return match
}

func IsWordValidator(input string) error {
	if input != "" && !match(isWordRegex, input) {
		return textinput.ErrInputValidation
	}
	return nil
}

func IsUsernameFolder(input string) error {
	if input != "" && !match(isUsernameFolder, input) {
		return textinput.ErrInputValidation
	}
	return nil
}

func IsUsernameFolderOrId(input string) error {
	if input != "" && (!match(isIdRegex, input) || !match(isUsernameFolder, input)) {
		return textinput.ErrInputValidation
	}
	return nil
}

func IsEmailValidator(input string) error {
	_, err := mail.ParseAddress(input)
	if err != nil {
		return textinput.ErrInputValidation
	}
	return nil
}

func IsPasswordValidator(input string) error {
	if input != "" && !match(isPasswordRegex, input) {
		return textinput.ErrInputValidation
	}
	return nil
}

func IsValidPathPart(input string) error {
	if input == "" || input == ".." || input == "." {
		return nil
	}
	if match(isWindowsPathBanChar, input) {
		return textinput.ErrInputValidation
	}
	if match(isLinuxPathBanChar, input) {
		return textinput.ErrInputValidation
	}
	if match(isWindowsPathBanName, input) {
		return textinput.ErrInputValidation
	}
	if match(isWindowsPathBanEnd, input) {
		return textinput.ErrInputValidation
	}
	return nil
}

func IsPathValidator(input string) error {
	parts := GetPathParts(input)
	if match(isWindowsDrive, parts[0]) {
		parts = parts[1:]
	}
	for _, p := range parts {
		if IsValidPathPart(p) != nil {
			return textinput.ErrInputValidation
		}
	}
	return nil
}

func IsGlobValidator(input string) error {
	_, err := filepath.Match(input, "")
	if err != nil {
		return textinput.ErrInputValidation
	}
	return nil
}

func IsMimeTypeValidator(input string) error {
	_, _, err := mime.ParseMediaType(input)
	if err != nil {
		return textinput.ErrInputValidation
	}
	return nil
}

func GetDataSizeValidator(max uint64) func(string) error {
	return func(input string) error {
		if input == "" {
			return nil
		}
		v, err := humanize.ParseBytes(input)
		if err != nil {
			return textinput.ErrInputValidation
		}
		if max != 0 && v > max {
			return textinput.ErrInputValidation
		}
		return nil
	}
}

func GetValidValueArrayValidator(allowEmpty bool, itemValidator func(string) error) func(string) error {
	return func(input string) error {
		if input == "" {
			if allowEmpty {
				return nil
			}
			return textinput.ErrInputValidation
		}
		var values []string
		for _, v := range separator.Split(input, -1) {
			if v == "" {
				continue
			}
			if itemValidator(v) != nil {
				return textinput.ErrInputValidation
			}
			values = append(values, v)
		}
		return nil
	}
}
