package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/erikgeiser/promptkit/textinput"
	"net/mail"
	"sherry/shr/api"
	"sherry/shr/config"
	"sherry/shr/helpers"
	"strings"
)

type UserCredentials = struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse = struct {
	Message    []string `json:"message"`
	StatusCode int      `json:"statusCode"`
}

type SuccessRegistrationResponse = struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

var isWordRegex = regexp2.MustCompile(`^\w+$`, 0).MatchString
var isPasswordRegex = regexp2.MustCompile(`(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}`, 0).MatchString

func match(regex func(s string) (bool, error), input string) bool {
	match, _ := regex(input)
	return match
}

func isWordValidator(input string) error {
	if input != "" && !match(isWordRegex, input) {
		return textinput.ErrInputValidation
	}
	return nil
}

func isEmailValidator(input string) error {
	_, err := mail.ParseAddress(input)
	if err != nil {
		return textinput.ErrInputValidation
	}
	return nil
}

func isValidPassword(input string) error {
	if input != "" && !match(isPasswordRegex, input) {
		return textinput.ErrInputValidation
	}
	return nil
}

func getUserInfo(register bool, email string, password string, user string) UserCredentials {
	if email == "" {
		input := textinput.New("Email")
		input.Validate = isEmailValidator
		email, _ = input.RunPrompt()
	}
	if register && user == "" {
		input := textinput.New("Username")
		input.Validate = isWordValidator
		user, _ = input.RunPrompt()
	}
	if password == "" {
		input := textinput.New("Password")
		input.Placeholder = "At least 6 letters long, one Capital letter, one lowercase letter, and one number"
		input.Validate = isValidPassword
		input.Hidden = true
		password, _ = input.RunPrompt()
	}
	return UserCredentials{Username: user, Email: email, Password: password}
}

func checkUserExists(email string, user string) bool {
	authConfig := config.GetAuthConfig()

	for _, u := range authConfig.Sources {
		if u.Email == email || u.Username == user {
			return true
		}
	}
	return false
}

func validateResponse(res string, err error) (string, error) {
	if err != nil {
		if res != "" {
			var resErr ErrorResponse
			err := json.Unmarshal([]byte(res), &resErr)
			if err != nil {
				helpers.PrintErr(res)
			}
			helpers.PrintErr(strings.Join(resErr.Message, "\n"))
		} else {
			helpers.PrintErr(err.Error())
		}
		return "", err
	}
	return res, nil
}

func RegisterUser(email string, password string, user string) bool {
	info := getUserInfo(true, email, password, user)

	if checkUserExists(info.Email, info.Username) {
		helpers.PrintErr("User already authorized")
	}

	helpers.PrintMessage("Creating user...")

	body, _ := json.Marshal(info)
	res, err := api.Post("auth/sign-up", body)
	if _, err := validateResponse(res, err); err != nil {
		return false
	}

	var createdUser SuccessRegistrationResponse
	err = json.Unmarshal([]byte(res), &createdUser)

	helpers.PrintMessage("User created successfully")

	if !LoginUser(createdUser.Email, info.Password) {
		return false
	}

	if config.GetAuthConfig().Default == "" {
		helpers.PrintMessage("It is the only user, setting it as default...")
		SetDefaultUser(info.Username)
	}

	return true
}

func LoginUser(email string, password string) bool {
	info := getUserInfo(false, email, password, "")

	helpers.PrintMessage("Authorizing...")

	body, _ := json.Marshal(map[string]string{"email": info.Email, "password": info.Password})
	res, err := api.Post("auth/sign-in", body)
	if _, err := validateResponse(res, err); err != nil {
		return false
	}

	var authResponse config.Credentials
	err = json.Unmarshal([]byte(res), &authResponse)

	authConfig := config.GetAuthConfig()
	authConfig.Sources[authResponse.Id] = authResponse
	return true
}

func SetDefaultUser(user string) bool {
	authConfig := config.GetAuthConfig()

	var credentials *config.Credentials = nil
	for _, u := range authConfig.Sources {
		if u.Username == user {
			credentials = &u
			break
		}
	}

	if credentials == nil {
		helpers.PrintErr("User not found")
		return false
	}

	if authConfig.Default == credentials.Id {
		helpers.PrintErr("User is already default")
		return false
	}

	authConfig.Default = credentials.Id

	helpers.PrintMessage(fmt.Sprintf("User %s set as default", credentials.Username))

	return true
}

func PrintDefaultUser() bool {
	authConfig := config.GetAuthConfig()
	if authConfig.Default == "" {
		helpers.PrintMessage("No default user set")
		return false
	}

	for _, u := range authConfig.Sources {
		if u.Id == authConfig.Default {
			helpers.PrintMessage(fmt.Sprintf("Default user: %s", u.Username))
			return false
		}
	}

	return false
}

func PrintUsers() bool {
	authConfig := config.GetAuthConfig()

	helpers.PrintMessage("* - default user\n")
	for _, u := range authConfig.Sources {
		isDefault := " "
		if u.Id == authConfig.Default {
			isDefault = "*"
		}
		helpers.PrintMessage(fmt.Sprintf("%s %s", isDefault, u.Username))
	}

	return false
}
