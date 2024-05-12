package auth

import (
	"encoding/json"
	"fmt"
	"sherry/shr/api"
	"sherry/shr/config"
	"sherry/shr/helpers"
)

type UserCredentials = struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SuccessRegistrationResponse = struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func getUserInfo(register bool, email string, password string, user string) UserCredentials {
	email = helpers.Input("Email", email, helpers.IsEmailValidator, "", false)
	if register {
		user = helpers.Input("Username", user, helpers.IsWordValidator, "", false)
	}
	password = helpers.Input(
		"Password",
		password,
		helpers.IsPasswordValidator,
		"At least 6 letters long, one Capital letter, one lowercase letter, and one number",
		true,
	)
	return UserCredentials{
		Email:    email,
		Username: user,
		Password: password,
	}
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

func getUserString(user config.Credentials) string {
	return fmt.Sprintf("%s(%s)", user.Username, user.Email)
}

func RegisterUser(email string, password string, user string) bool {
	info := getUserInfo(true, email, password, user)

	if checkUserExists(info.Email, info.Username) {
		helpers.PrintErr("User already authorized")
	}

	helpers.PrintMessage("Creating user...")

	body, _ := json.Marshal(info)
	res, err := api.Post("auth/sign-up", body, "")
	if _, err := helpers.ValidateResponse(res, err); err != nil {
		return false
	}

	var createdUser SuccessRegistrationResponse
	err = json.Unmarshal([]byte(res), &createdUser)

	helpers.PrintMessage("User created successfully")

	return LoginUser(createdUser.Email, info.Password)
}

func LoginUser(email string, password string) bool {
	info := getUserInfo(false, email, password, "")

	helpers.PrintMessage("Authorizing...")

	body, _ := json.Marshal(map[string]string{"email": info.Email, "password": info.Password})
	res, err := api.Post("auth/sign-in", body, "")
	if _, err := helpers.ValidateResponse(res, err); err != nil {
		return false
	}

	var authResponse config.Credentials
	err = json.Unmarshal([]byte(res), &authResponse)

	authConfig := config.GetAuthConfig()
	authConfig.Sources[authResponse.UserId] = authResponse

	if authConfig.Default == "" {
		helpers.PrintMessage("It is the only user, setting it as default...")
		SetDefaultUser(authResponse.Username)
	}

	helpers.PrintMessage("User was successfully logged in")

	return true
}

func FindUserByUsername(username string, withDefault bool) *config.Credentials {
	authConfig := config.GetAuthConfig()
	if username == "" && withDefault {
		username = authConfig.Sources[authConfig.Default].Username
	}
	if username == "" {
		return nil
	}
	for _, u := range authConfig.Sources {
		if u.Username == username {
			return &u
		}
	}
	return nil
}

func GetUserById(userId string) *config.Credentials {
	if v, ok := config.GetAuthConfig().Sources[userId]; ok {
		return &v
	}
	return nil
}

func SetDefaultUser(user string) bool {
	authConfig := config.GetAuthConfig()

	var credentials = FindUserByUsername(user, false)

	if credentials == nil {
		helpers.PrintErr("User not found")
		return false
	}

	if authConfig.Default == credentials.UserId {
		helpers.PrintErr("User is already default")
		return false
	}

	authConfig.Default = credentials.UserId

	helpers.PrintMessage(fmt.Sprintf("User %s set as default", getUserString(*credentials)))

	return true
}

func PrintDefaultUser() bool {
	authConfig := config.GetAuthConfig()
	if authConfig.Default == "" {
		helpers.PrintMessage("No default user set")
		return false
	}

	for _, u := range authConfig.Sources {
		if u.UserId == authConfig.Default {
			helpers.PrintMessage(fmt.Sprintf("Default user: %s", getUserString(u)))
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
		if u.UserId == authConfig.Default {
			isDefault = "*"
		}
		helpers.PrintMessage(fmt.Sprintf("%s %s", isDefault, getUserString(u)))
	}

	return false
}
