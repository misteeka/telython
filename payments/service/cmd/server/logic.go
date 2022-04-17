package main

import (
	"fmt"
	"strconv"
	"telython/payments/service/pkg/accounts"
	"telython/payments/service/pkg/database"
	"telython/payments/service/pkg/log"
	"telython/payments/service/pkg/payments"
	"telython/payments/service/pkg/status"
	"telython/pkg/eplidr"
)

func getUsername(accountId uint64) (string, bool, error) {
	return database.Accounts.GetString(eplidr.Column{Key: "id", Value: accountId}, "name")
}

func sendPayment(senderId uint64, receiverId uint64, amount uint64, timestamp uint64) status.Status {
	// check amount
	if amount <= 0 {
		return status.WRONG_AMOUNT
	}

	// check currency code mismatch
	sender, err := accounts.Load(senderId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR
	}
	if sender == nil {
		return status.NOT_FOUND
	}
	receiver, err := accounts.Load(receiverId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR
	}
	if receiver == nil {
		return status.NOT_FOUND
	}

	if sender.Currency != receiver.Currency {
		return status.CURRENCY_CODE_MISMATCH
	}

	balance, err := sender.GetBalance()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR
	}
	if balance < amount {
		return status.INSUFFICIENT_FUNDS
	}
	payment := payments.New(sender, receiver, amount, timestamp)

	err = payment.Commit()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR
	}
	go func() {
		err := database.Balances.Set(senderId, []string{"onSerial", "balance"}, []interface{}{payment.Timestamp, balance - amount})
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return
		}
		err = database.Balances.Set(receiverId, []string{"onSerial", "balance"}, []interface{}{payment.Timestamp, balance + amount})
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return
		}
	}()

	return status.SUCCESS
}

func getBalance(accountId uint64) (status.Status, uint64) {
	account, err := accounts.Load(accountId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR, 0
	}
	if account == nil {
		return status.NOT_FOUND, 0
	}
	balance, err := account.GetBalance()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR, 0
	}
	return status.SUCCESS, balance
}

func getHistory(accountId uint64) (status.Status, *[]payments.Payment) {
	var history []payments.Payment
	rows, err := database.Payments.Query(fmt.Sprintf("SELECT * FROM {table} WHERE `sender` = %d OR `receiver` = %d LIMIT 2000;", accountId, accountId), accountId)
	if err != nil {
		return status.INTERNAL_SERVER_ERROR, nil
	}
	for rows.Next() {
		var payment payments.Payment
		err = rows.Scan(&payment.Id, &payment.Sender, &payment.Receiver, &payment.Amount, &payment.Timestamp, &payment.Currency)
		if err != nil {
			return status.INTERNAL_SERVER_ERROR, nil
		}
		history = append(history, payment)
	}
	return status.SUCCESS, &history
}

func getAccountInfo(accountId uint64) (status.Status, *accounts.Account) {
	account, err := accounts.Load(accountId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR, nil
	}
	return status.SUCCESS, account
}

func getPayment(id uint64, accountId uint64) (status.Status, *payments.Payment) {
	payment := payments.Payment{
		Id: id,
	}
	err, found := database.Payments.Get(
		accountId,
		eplidr.PlainToColumns([]string{"id"}, []interface{}{id}),
		[]string{"sender", "receiver", "amount", "timestamp", "currency"},
		[]interface{}{&payment.Sender, &payment.Receiver, &payment.Amount, &payment.Timestamp, &payment.Currency},
	)
	if err != nil {
		return status.INTERNAL_SERVER_ERROR, nil
	}
	if !found {
		return status.NOT_FOUND, nil
	}
	return status.SUCCESS, &payment
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

func createAccount(username string, currency int, timestamp uint64) (status.Status, uint64) {
	accountId := fnv64(username + strconv.FormatUint(timestamp, 10))
	err := database.Accounts.Put(accountId, eplidr.PlainToColumns([]string{"id", "name", "currency"}, []interface{}{accountId, username, currency}))
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR, 0
	}
	err = database.Balances.Put(accountId, []string{"id", "balance"}, []interface{}{accountId, 1000000})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return status.INTERNAL_SERVER_ERROR, 0
	}
	return status.SUCCESS, accountId
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
