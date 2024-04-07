package auth

import (
	"github.com/erikgeiser/promptkit/textinput"
	"regexp"
	"sherry/shr/helpers"
)

type Credentials struct {
	Email        string `json:"email"`
	Nickname     string `json:"nickname"`
	RefreshToken string `json:"refreshToken"`
}

type AuthorizationConfig struct {
	Sources map[string]Credentials `json:"records"`
	Default string                 `json:"default"`
}

type UserCredentials = struct {
	User     string
	Password string
}

var isWord = regexp.MustCompile(`^\w+$`).MatchString

func isWordValidator(input string) error {
	if input != "" && !isWord(input) {
		return textinput.ErrInputValidation
	}
	return nil
}

func getUserInfo(user string, password string) UserCredentials {
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
	return UserCredentials{user, password}
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
