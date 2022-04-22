package client

import (
	"fmt"
	"telython/pkg/http"
	httpclient "telython/pkg/http/client"
)

var client *httpclient.Client

func init() {
	client = httpclient.New("127.0.0.1:8001", "/auth/")
}

func SignIn(username string, password string) (*http.Error, error) {
	json, err := client.Put("signIn", fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password))
	if err != nil {
		return nil, err
	}
	return httpclient.GetError(json), nil
}
func CheckPassword(username string, password string) (*http.Error, error) {
	json, err := client.Get(fmt.Sprintf("checkPassword?u=%s&p=%s", username, password))
	if err != nil {
		return nil, err
	}
	return httpclient.GetError(json), nil
}
func ResetPassword(username string, oldPassword string, newPassword string) (*http.Error, error) {
	json, err := client.Put("resetPassword", fmt.Sprintf(`{"username":"%s", "oldPassword":"%s", "newPassword":"%s"}`, username, oldPassword, newPassword))
	if err != nil {
		return nil, err
	}
	return httpclient.GetError(json), nil
}
func RequestSignUpCode(username string, email string) (*http.Error, error) {
	json, err := client.Post("requestSignUpCode", fmt.Sprintf(`{"username":"%s", "email":"%s"}`, username, email))
	if err != nil {
		return nil, err
	}
	return httpclient.GetError(json), nil
}
func RequestPasswordRecovery(username string) (*http.Error, error) {
	json, err := client.Put("requestPasswordRecovery", fmt.Sprintf(`{"username":"%s"}`, username))
	if err != nil {
		return nil, err
	}
	return httpclient.GetError(json), nil
}
func RecoverPassword(username string, newPassword string, code string) (*http.Error, error) {
	json, err := client.Put("recoverPassword", fmt.Sprintf(`{"username":"%s", "newPassword":"%s", "code":"%s"}`, username, newPassword, code))
	if err != nil {
		return nil, err
	}
	return httpclient.GetError(json), nil
}
func SignUp(username string, password string, code string) (*http.Error, error) {
	json, err := client.Post("signUp", fmt.Sprintf(`{"username":"%s", "password":"%s", "code":"%s"}`, username, password, code))
	if err != nil {
		return nil, err
	}
	return httpclient.GetError(json), nil
}
