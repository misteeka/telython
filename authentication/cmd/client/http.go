package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"telython/authentication/pkg/client"
	"telython/pkg/http"
)

func print(error *http.Error, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(http.ToReadable(error))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Telython auth")
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
			print(client.SignIn(username, password))
		} else if strings.Compare("checkPassword", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			print(client.CheckPassword(username, password))
		} else if strings.Compare("resetPassword", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			oldPassword := args[1]
			newPassword := args[2]
			print(client.ResetPassword(username, oldPassword, newPassword))
		} else if strings.Compare("requestSignUpCode", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			email := args[1]
			print(client.RequestSignUpCode(username, email))
		} else if strings.Compare("requestPasswordRecovery", cmd) == 0 {
			if len(args) < 1 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			print(client.RequestPasswordRecovery(username))
		} else if strings.Compare("recoverPassword", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			newPassword := args[1]
			code := args[2]
			print(client.RecoverPassword(username, newPassword, code))
		} else if strings.Compare("signUp", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			code := args[2]
			print(client.SignUp(username, password, code))
		} else {
			fmt.Println("Unknown command.")
		}
	}

}
