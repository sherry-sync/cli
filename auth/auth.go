package auth

import (
	"github.com/erikgeiser/promptkit/textinput"
	"regexp"
	"sherry/cli/helpers"
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

func getUserInfo(user string, password string) UserCredentials {
	input := textinput.New("Local profile name")
	input.Placeholder = "default"
	input.Validate = isWordValidator
	input.InitialValue = "default"
	profileName, _ := input.RunPrompt()

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

func RegisterUser(user string, password string) {
	info := getUserInfo(user, password)
	helpers.PrintJson(info)
}

func LoginUser(user string, password string) {
	info := getUserInfo(user, password)
	helpers.PrintJson(info)
}

func SetDefaultUser(user string) {

}
