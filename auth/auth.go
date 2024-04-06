package auth

import (
	"github.com/erikgeiser/promptkit/textinput"
	"regexp"
	"sherry/shr/helpers"
)

type UserCredentials = struct {
	LocalProfile string
	User         string
	Password     string
}

var isWord = regexp.MustCompile(`^\w+$`).MatchString

func isWordValidator(input string) error {
	if input != "" && !isWord(input) {
		return textinput.ErrInputValidation
	}
	return nil
}

func getUserInfo(profileName string, user string, password string) UserCredentials {
	if profileName == "" {
		input := textinput.New("Local profile name")
		input.Placeholder = "default"
		input.Validate = isWordValidator
		input.InitialValue = "default"
		profileName, _ = input.RunPrompt()
	}
	if user == "" {
		input := textinput.New("Username")
		input.Validate = isWordValidator
		user, _ = input.RunPrompt()
	}
	if password == "" {
		input := textinput.New("Password")
		input.Hidden = true
		password, _ = input.RunPrompt()
	}
	return UserCredentials{profileName, user, password}
}

func RegisterUser(profileName string, user string, password string) {
	info := getUserInfo(profileName, user, password)
	helpers.PrintJson(info)
}

func LoginUser(profileName string, user string, password string) {
	info := getUserInfo(profileName, user, password)
	helpers.PrintJson(info)
}

func SetDefaultUser(user string) {

}
