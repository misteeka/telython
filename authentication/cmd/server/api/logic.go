package api

import (
	"fmt"
	"math/rand"
	"strconv"
	"telython/authentication/cmd/server/mail"
	"telython/authentication/pkg/database"
	"telython/pkg/eplidr"
	"telython/pkg/http"
	"telython/pkg/log"
	"time"
)

func requestEmailCode(username string) *http.Error {
	email, found, err := database.UsersByName.GetString(username, "email")
	if err != nil {
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return http.ToError(http.NOT_FOUND)
	}
	fmt.Println(email)
	// sendEmail(email, "", "")
	return nil
}
func isEmailExists(email string) (bool, error) {
	_, found, err := database.UsersByEmail.GetString(email, "name")
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	} else {
		return true, nil
	}
}
func isUsernameExists(username string) (bool, error) {
	_, found, err := database.UsersByName.GetString(username, "password")
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	} else {
		return true, nil
	}
}

func signIn(username string, password string, ip string) *http.Error {
	checkResponse := checkPassword(username, password)
	if checkResponse == nil {
		err := database.UsersByName.Set(username, eplidr.Columns{{"last_ip", ip}, {"last_login", time.Now().UnixMicro()}})
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR)
		}
		return nil
	} else {
		return checkResponse
	}
}
func checkPassword(username string, password string) *http.Error {
	val, found, err := database.UsersByName.GetString(username, "password")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Account Not Found!",
		}
	}
	if val != password {
		return &http.Error{
			Code:    http.AUTHORIZATION_FAILED,
			Message: "Wrong Password!",
		}
	}
	return nil
}
func resetPassword(username string, oldPassword string, newPassword string) *http.Error {
	checkResponse := checkPassword(username, oldPassword)
	if checkResponse == nil {
		err := database.UsersByName.SingleSet(username, "password", newPassword)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR)
		}
		return nil
	} else {
		return checkResponse
	}
}
func requestSignUpCode(username string, email string, ip string) *http.Error {
	exists, err := isEmailExists(email)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if exists {
		return &http.Error{
			Code:    http.ALREADY_EXISTS,
			Message: "This Email Is Already Registered!",
		}
	}
	code := rand.Intn(999999)
	err = database.PendingEmailConfirmations.Put(username, []string{"name", "email", "code", "timestamp"}, []interface{}{username, email, code, time.Now().UnixMicro()})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	go func() {
		err := mail.Send(email, "Confirm email", fmt.Sprintf(`<h2>Telython registration</h2>
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
	return nil
}
func requestPasswordRecovery(username string) *http.Error {
	code := strconv.FormatInt(int64(rand.Intn(999999)), 10)
	email, found, err := database.UsersByName.GetString(username, "email")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Account Not Found!",
		}
	}
	err = database.EmailCodes.Put(username, []string{"code"}, []interface{}{code})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	go func() {
		err := mail.Send(email, "Code from Telython", fmt.Sprintf(`
		<h2>Email confirmation</h2>
		<div>Use code <b>%s</b> to confirm the email for password recovery.</div>
		<div>Enter this code in the registration form in your app.	</div>
		<div>If you did not request a password recovery, please ignore this message.</div>
	`, code))
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	return nil
}
func recoverPassword(username string, code string, newPassword string) *http.Error {
	savedCode, found, err := database.EmailCodes.GetString(username, "code")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Code Not Found!",
		}
	}
	if savedCode == code {
		err := database.UsersByEmail.SingleSet(username, "password", newPassword)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR)
		}
		return nil
	} else {
		return &http.Error{
			Code:    http.AUTHORIZATION_FAILED,
			Message: "Entered Code Is Wrong!",
		}
	}
}
func signUp(username string, password string, code string, ip string) *http.Error {
	exists, err := isUsernameExists(username)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if exists {
		return &http.Error{
			Code:    http.ALREADY_EXISTS,
			Message: "This Username Is Already Registered!",
		}
	}
	email, found, err := database.PendingEmailConfirmations.GetString(username, "email")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return &http.Error{
			Code:    http.ALREADY_EXISTS,
			Message: "This Email Is Already Registered!",
		}
	}
	exists, err = isEmailExists(email)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if exists {
		return http.ToError(http.ALREADY_EXISTS)
	}

	savedCode, found, err := database.PendingEmailConfirmations.GetString(username, "code")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return http.ToError(http.NOT_FOUND)
	}
	if code == savedCode {
		err = database.UsersByName.Put(username, []string{"name", "password", "email", "reg_ip", "last_ip", "reg_date", "last_login"}, []interface{}{username, password, email, ip, ip, time.Now().UnixMicro(), time.Now().UnixMicro()})
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR)
		}
		err = database.UsersByEmail.Put(email, []string{"name", "email"}, []interface{}{username, email})
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR)
		}
		err := database.PendingEmailConfirmations.Remove(username)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return nil
		}
		go func() {
			err = mail.Send(email, "You was register on Telython!", fmt.Sprintf(`<h2>Telython registration</h2>
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
		return &http.Error{
			Code:    http.AUTHORIZATION_FAILED,
			Message: "Entered Code Is Wrong!",
		}
	}
	return nil
}
