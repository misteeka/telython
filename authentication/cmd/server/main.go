package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"telython/authentication/pkg/cfg"
	"telython/authentication/pkg/database"
	"telython/authentication/pkg/log"
	"telython/authentication/pkg/status"
	"telython/pkg/eplidr"
	"time"
)

func requestEmailCode(username string) status.Status {
	email, found, err := database.SingleUsers.GetString(username, "email")
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return status.NOT_FOUND
	}
	fmt.Println(email)
	// sendEmail(email, "", "")
	return status.SUCCESS
}
func isEmailExists(email string) status.Status {
	_, found, err := database.Users.GetString(eplidr.Column{Key: "email", Value: email}, "password")
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return status.NO
	} else {
		return status.YES
	}
}
func isUsernameExists(username string) status.Status {
	_, found, err := database.SingleUsers.GetString(username, "password")
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return status.NO
	} else {
		return status.YES
	}
}

func signIn(username string, password string, ip string) status.Status {
	checkResponse := checkPassword(username, password)
	if checkResponse == status.SUCCESS {
		err := database.SingleUsers.Set(username, []string{"last_ip", "last_login"}, []interface{}{ip, time.Now().UnixMicro()})
		if err != nil {
			return status.INTERNAL_SERVER_ERROR
		}
		return status.SUCCESS
	} else {
		return checkResponse
	}
}
func checkPassword(username string, password string) status.Status {
	val, found, err := database.SingleUsers.GetString(username, "password")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return status.AUTHORIZATION_FAILED
	}
	if val != password {
		return status.AUTHORIZATION_FAILED
	}
	return status.SUCCESS
}
func resetPassword(username string, oldPassword string, newPassword string) status.Status {
	checkResponse := checkPassword(username, oldPassword)
	if checkResponse == status.SUCCESS {
		err := database.SingleUsers.SingleSet(username, "password", newPassword)
		if err != nil {
			return status.INTERNAL_SERVER_ERROR
		}
		return status.SUCCESS
	} else {
		return checkResponse
	}
}
func requestSignUpCode(username string, email string, ip string) status.Status {
	isEmailExists := isEmailExists(email)
	if isEmailExists == status.YES {
		return status.ALREADY_EXISTS
	} else if isEmailExists == status.INTERNAL_SERVER_ERROR {
		return status.INTERNAL_SERVER_ERROR
	}
	code := rand.Intn(999999)
	go func() {
		err := sendEmail(email, "Confirm email", fmt.Sprintf(`<h2>Telython registration</h2>
	<div>
		<div>Hello, %s.</div>
		<div>Use code <b>%d</b> to confirm the email for registration.</div>
		<div>Enter this code in the registration form in your app.</div>
		<div>If you did not request a registration, please ignore this message.</div>
	</div>`, username, code))
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	err := database.PendingEmailConfirmations.Put(username, []string{"email", "code", "timestamp"}, []interface{}{email, code, time.Now().UnixMicro()})
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	return status.SUCCESS
}
func requestPasswordRecovery(username string) status.Status {
	code := strconv.FormatInt(int64(rand.Intn(999999)), 10)
	email, found, err := database.SingleUsers.GetString(username, "email")
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return status.NOT_FOUND
	}
	err = database.EmailCodes.Put(username, []string{"code"}, []interface{}{code})
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	go func() {
		err := sendEmail(email, "Code from Telython", fmt.Sprintf(`
		<h2>Email confirmation</h2>
		<div>Use code <b>%s</b> to confirm the email for password recovery.</div>
		<div>Enter this code in the registration form in your app.	</div>
		<div>If you did not request a password recovery, please ignore this message.</div>
	`, code))
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	return status.SUCCESS
}
func recoverPassword(username string, code string, newPassword string) status.Status {
	savedCode, found, err := database.EmailCodes.GetString(username, "code")
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return status.NOT_FOUND
	}
	if savedCode == code {
		err := database.SingleUsers.SingleSet(username, "password", newPassword)
		if err != nil {
			return status.INTERNAL_SERVER_ERROR
		}
		return status.SUCCESS
	} else {
		return status.AUTHORIZATION_FAILED
	}
}
func signUp(username string, password string, code string, ip string) status.Status {
	isUsernameExists := isUsernameExists(username)
	if isUsernameExists == status.YES {
		return status.ALREADY_EXISTS
	} else if isUsernameExists == status.INTERNAL_SERVER_ERROR {
		return status.INTERNAL_SERVER_ERROR
	}
	email, found, err := database.PendingEmailConfirmations.GetString(username, "email")
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return status.NOT_FOUND
	}
	isEmailExists := isEmailExists(email)
	if isEmailExists == status.YES {
		return status.ALREADY_EXISTS
	} else if isEmailExists == status.INTERNAL_SERVER_ERROR {
		return status.INTERNAL_SERVER_ERROR
	}
	savedCode, found, err := database.PendingEmailConfirmations.GetString(username, "code")
	if err != nil {
		return status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return status.NOT_FOUND
	}
	if code == savedCode {
		err = database.SingleUsers.Put(username, []string{"password", "email", "reg_ip", "last_ip", "reg_date", "last_login"}, []interface{}{password, email, ip, ip, time.Now().UnixMicro(), time.Now().UnixMicro()})
		if err != nil {
			return status.INTERNAL_SERVER_ERROR
		}
		database.PendingEmailConfirmations.Remove(username)
		go func() {
			err = sendEmail(email, "You was register on Telython!", fmt.Sprintf(`<h2>Telython registration</h2>
		<div>
			<div>Hello, %s.</div>
			<div>Your was registered on Telython.</div>
			<div>Enjoy messaging!</div>
		</div>`, username))
			if err != nil {
				log.ErrorLogger.Println(err.Error())
			}
		}()
	} else {
		return status.AUTHORIZATION_FAILED
	}
	return status.SUCCESS
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.GOMAXPROCS(8)
	rand.Seed(time.Now().UnixNano())

	log.InfoLogger.Println("Starting...")
	log.InfoLogger.Println("Config loading")
	panicIfError(cfg.LoadConfig())
	log.InfoLogger.Println("Database start")
	panicIfError(database.InitDatabase())

	log.InfoLogger.Println("Gomail start")
	initMailClient()
	// TestMail()

	log.InfoLogger.Println("TLS initialization")
	initTLS()

	panicIfError(initFiber())
	ListenAndServe() // Blocking

	log.InfoLogger.Println("Shutdown...")
	log.InfoLogger.Println("Goodbye!")
}
