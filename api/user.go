package api

import (
	"encoding/json"
	"fmt"
)

type ResponseUser = struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type ResponseLogin = struct {
	UserId       string `json:"userId"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type PayloadUser = struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PayloadLogin = struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserRegister(payload PayloadUser) (*ResponseUser, error) {
	body, _ := json.Marshal(payload)
	res, err := ValidateResponse(Post("auth/sign-up", body, ""))
	if err != nil {
		return nil, err
	}

	return ParseResponse[ResponseUser](res)
}

func UserLogin(payload PayloadLogin) (*ResponseLogin, error) {
	body, _ := json.Marshal(payload)
	res, err := ValidateResponse(Post("auth/sign-in", body, ""))
	if err != nil {
		return nil, err
	}

	return ParseResponse[ResponseLogin](res)
}

func UserGet(accessToken string) (*ResponseUser, error) {
	res, err := ValidateResponse(Get("/user/me", accessToken))
	if err != nil {
		return nil, err
	}

	return ParseResponse[ResponseUser](res)
}

func UserUpdate(payload PayloadUser, accessToken string) (*ResponseUser, error) {
	body, _ := json.Marshal(payload)
	res, err := ValidateResponse(Patch("/user/me", body, accessToken))
	if err != nil {
		return nil, err
	}

	return ParseResponse[ResponseUser](res)
}

func UserFind(username string, accessToken string) (*ResponseUser, error) {
	res, err := ValidateResponse(Get(fmt.Sprintf("/user/find?username=%s", username), accessToken))
	if err != nil {
		return nil, err
	}

	return ParseResponse[ResponseUser](res)
}
