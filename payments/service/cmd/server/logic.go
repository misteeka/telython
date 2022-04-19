package main

import (
	"fmt"
	"math/big"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/accounts"
	"telython/payments/service/pkg/database"
	"telython/payments/service/pkg/payments"
	"telython/pkg/eplidr"
	"telython/pkg/http"
	"telython/pkg/log"
)

func getUsername(nameHash uint64) (string, bool, error) {
	return database.Accounts.GetString(nameHash, "name")
}

func sendPayment(senderName string, receiverName string, currency *currency.Currency, timestamp uint64) *http.Error {
	// check amount
	if currency.Amount.Cmp(big.NewInt(0)) <= 0 {
		return &http.Error{
			Code:    http.WRONG_AMOUNT,
			Message: "Amount Must Be More Than 0 ",
		}
	}

	// check currency code mismatch
	senderId, found, err := accounts.GetId(senderName)
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
	receiverId, found, err := accounts.GetId(receiverName)
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

	balance, err := accounts.GetBalance(senderId, currency.Type.Id)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if balance.Cmp(currency.Amount) < 0 {
		return &http.Error{
			Code:    http.INSUFFICIENT_FUNDS,
			Message: "Insufficient Funds, Top Up Your Balance First",
		}
	}
	payment := payments.New(senderId, receiverId, currency, timestamp)

	err = payment.Commit()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return nil
}

func addPayment(receiverName string, currency *currency.Currency, timestamp uint64) *http.Error {
	// check amount
	if currency.Amount.Cmp(big.NewInt(0)) <= 0 {
		return &http.Error{
			Code:    http.WRONG_AMOUNT,
			Message: "Amount Must Be More Than 0 ",
		}
	}

	receiverId, found, err := accounts.GetId(receiverName)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Receiver Not Found!",
		}
	}

	payment := payments.New(0, receiverId, currency, timestamp)

	err = payment.Commit()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return nil
}

func getBalance(username string, currencyCode uint64) (*http.Error, *big.Int) {
	accountId, found, err := accounts.GetId(username)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	if !found {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Account Not Found!",
		}, nil
	}
	balance, err := accounts.GetBalance(accountId, currencyCode)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	return nil, balance
}

func getHistory(username string) (*http.Error, *[]payments.Payment) {
	accountId, found, err := accounts.GetId(username)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	if !found {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Account Not Found!",
		}, nil
	}
	var history []payments.Payment
	rows, err := database.Payments.Query(fmt.Sprintf("SELECT * FROM {table} WHERE `sender` = %d OR `receiver` = %d LIMIT 2000;", accountId, accountId), accountId)
	if err != nil {
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	for rows.Next() {
		var payment payments.Payment
		var serializedAmount string
		var currencyCode uint64
		err = rows.Scan(&payment.Id, &payment.Sender, &payment.Receiver, &serializedAmount, &currencyCode, &payment.Timestamp)
		if err != nil {
			rows.Close()
			return http.ToError(http.INTERNAL_SERVER_ERROR), nil
		}
		currency, err := currency.Deserialize(currencyCode, serializedAmount)
		if err != nil {
			rows.Close()
			return http.ToError(http.INTERNAL_SERVER_ERROR), nil
		}
		payment.Currency = currency
		history = append(history, payment)
	}
	return nil, &history
}

func getAccountInfo(username string) (*http.Error, *accounts.AccountInfo) {
	accountInfo, err := accounts.GetAccountInfo(fnv64(username))
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	return nil, accountInfo
}

func getPayment(id uint64, username string) (*http.Error, *payments.Payment) {
	payment := payments.Payment{
		Id: id,
	}
	var currencyCode uint64
	var amountBase64 string
	nameHash := fnv64(username)
	err, found := database.Payments.Get(
		nameHash,
		eplidr.PlainToColumns([]string{"id"}, []interface{}{id}),
		[]string{"sender", "receiver", "amount", "timestamp", "currency"},
		[]interface{}{&payment.Sender, &payment.Receiver, &amountBase64, &payment.Timestamp, &currencyCode},
	)
	if err != nil {
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	if !found {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Payment Not Found!",
		}, nil
	}

	return nil, &payment
}

func fnv64(key string) uint64 {
	hash := uint64(4332272522)
	const prime64 = uint64(33555238)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime64
		hash ^= uint64(key[i])
	}
	return hash
}

func createAccount(username string) *http.Error {
	nameHash := fnv64(username)
	_, found, err := database.Accounts.GetString(nameHash, "name")
	if err != nil {
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if found {
		return &http.Error{
			Code:    http.ALREADY_EXISTS,
			Message: "Account For User Already Created!",
		}
	}
	err = database.Accounts.Put(nameHash, []string{"name", "nameHash"}, []interface{}{username, nameHash})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return nil
}

/*
tx, err := database.Accounts.RawTx(senderId)
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR
	}

	sender, err := account.Load(senderId, tx)
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR
	}
	if sender == nil {
		tx.Rollback()
		return status.NOT_FOUND
	}
	receiver, err := account.Load(receiverId, tx)
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR
	}
	if receiver == nil {
		tx.Rollback()
		return status.NOT_FOUND
	}

	if sender.Currency != receiver.Currency {
		return status.CURRENCY_CODE_MISMATCH
	}
	if sender.Balance < amount {
		return status.INSUFFICIENT_FUNDS
	}
	payment := payment.New(sender, receiver, amount, tx, timestamp)
	err = payment.Transfer()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		payment.Fail()
		return status.INTERNAL_SERVER_ERROR
	}
	err = payment.Commit()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		payment.Fail()
		return status.INTERNAL_SERVER_ERROR
	}
	return status.SUCCESS
*/
