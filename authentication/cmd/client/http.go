package main

import (
	"bufio"
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
	"strings"
	"time"
)

type Status int

var (
	SUCCESS               Status = 100
	INVALID_REQUEST       Status = 101
	INTERNAL_SERVER_ERROR Status = 102
	AUTHORIZATION_FAILED  Status = 103
	ALREADY_EXISTS        Status = 104
	NOT_FOUND             Status = 105
)

func statusToString(status Status) string {
	if status == SUCCESS {
		return "SUCCESS"
	}
	if status == INVALID_REQUEST {
		return "INVALID REQUEST"
	}
	if status == INTERNAL_SERVER_ERROR {
		return "INTERNAL SERVER ERROR"
	}
	if status == AUTHORIZATION_FAILED {
		return "AUTHORIZATION FAILED"
	}
	if status == ALREADY_EXISTS {
		return "ALREADY EXISTS"
	}
	if status == NOT_FOUND {
		return "NOT FOUND"
	}
	return fmt.Sprintf("%d", status)
}

var client fasthttp.HostClient

func init() {
	client = fasthttp.HostClient{
		Addr:                "127.0.0.1:8001",
		MaxIdleConnDuration: time.Minute,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
	}
}

func get(function string) (Status, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8001/auth/" + function)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return 0, err
	}
	response := resp.Body()
	ReleaseResponse(resp)
	i, err := strconv.Atoi(string(response))
	if err != nil {
		return 0, err
	}
	return Status(i), nil
}
func post(function string, json string) (Status, error) {
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte(json))
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI("http://127.0.0.1:8001/auth/" + function)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	if err != nil {
		return 0, err
	}
	response := resp.Body()
	ReleaseResponse(resp)
	i, err := strconv.Atoi(string(response))
	if err != nil {
		return 0, err
	}
	return Status(i), nil
}
func put(function string, json string) (Status, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8001/auth/" + function)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes([]byte("PUT"))
	req.SetBody([]byte(json))
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return 0, err
	}
	response := resp.Body()
	ReleaseResponse(resp)
	i, err := strconv.Atoi(string(response))
	if err != nil {
		return 0, err
	}
	return Status(i), nil
}
func delete(function string, json string) (Status, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8001/auth/" + function)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes([]byte("DELETE"))
	req.SetBody([]byte(json))
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	if err != nil {
		return 0, err
	}
	response := resp.Body()
	ReleaseResponse(resp)
	i, err := strconv.Atoi(string(response))
	if err != nil {
		return 0, err
	}
	return Status(i), nil
}
func ReleaseResponse(response *fasthttp.Response) {
	fasthttp.ReleaseResponse(response)
}

func SignIn(username string, password string) (Status, error) {
	return put("signIn", fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password))
}
func CheckPassword(username string, password string) (Status, error) {
	return get(fmt.Sprintf("checkPassword?u=%s&p=%s", username, password))
}
func ResetPassword(username string, oldPassword string, newPassword string) (Status, error) {
	return put("resetPassword", fmt.Sprintf(`{"username":"%s", "oldPassword":"%s", "newPassword":"%s"}`, username, oldPassword, newPassword))
}
func RequestSignUpCode(username string, email string) (Status, error) {
	return post("requestSignUpCode", fmt.Sprintf(`{"username":"%s", "email":"%s"}`, username, email))
}
func RequestPasswordRecovery(username string) (Status, error) {
	return put("requestPasswordRecovery", fmt.Sprintf(`{"username":"%s"}`, username))
}
func RecoverPassword(username string, newPassword string, code string) (Status, error) {
	return put("recoverPassword", fmt.Sprintf(`{"username":"%s", "newPassword":"%s", "code":"%s"}`, username, newPassword, code))
}
func SignUp(username string, password string, code string) (Status, error) {
	return post("signUp", fmt.Sprintf(`{"username":"%s", "password":"%s", "code":"%s"}`, username, password, code))
}

func print(status Status, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(statusToString(status))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Telyauth Shell")
	fmt.Println("---------------------")
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ReplaceAll(text, "\r", "")
		parts := strings.Split(text, " ")
		if len(parts) < 1 {
			fmt.Println("Wrong command")
			continue
		}
		cmd := parts[0]
		args := parts[1:]
		if strings.Compare("signIn", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			print(SignIn(username, password))
		} else if strings.Compare("checkPassword", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			print(CheckPassword(username, password))
		} else if strings.Compare("resetPassword", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			oldPassword := args[1]
			newPassword := args[2]
			print(ResetPassword(username, oldPassword, newPassword))
		} else if strings.Compare("requestSignUpCode", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			email := args[1]
			print(RequestSignUpCode(username, email))
		} else if strings.Compare("requestPasswordRecovery", cmd) == 0 {
			if len(args) < 1 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			print(RequestPasswordRecovery(username))
		} else if strings.Compare("recoverPassword", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			newPassword := args[1]
			code := args[2]
			print(RecoverPassword(username, newPassword, code))
		} else if strings.Compare("signUp", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			code := args[2]
			print(SignUp(username, password, code))
		} else {
			fmt.Println("Unknown command.")
		}
	}

}
