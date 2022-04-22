package main

import (
	"math/big"
	exchange "telython/payments/exchange/pkg/client"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/accounts"
	"telython/payments/service/pkg/database"
	"telython/payments/service/pkg/payments"
	"telython/pkg/cfg"
	"telython/pkg/eplidr"
	"telython/pkg/http"
	"telython/pkg/log"
	"telython/pkg/utils"
)

func getUsername(id uint64) (string, bool, error) {
	return database.Accounts.GetString(id, "name")
}

func sendPayment(senderId uint64, receiverId uint64, currencyFrom *currency.Currency, currencyCodeTo uint64, timestamp uint64) *http.Error {
	// Check amount
	if currencyFrom.Amount.Cmp(big.NewInt(0)) <= 0 {
		return &http.Error{
			Code:    http.WRONG_AMOUNT,
			Message: "Amount Must Be More Than 0 ",
		}
	}

	// Check do accounts exists
	exists, err := accounts.Exists(senderId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !exists {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Sender Not Found!",
		}
	}
	exists, err = accounts.Exists(receiverId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !exists {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Receiver Not Found!",
		}
	}

	var currencyTo *currency.Currency
	if currencyFrom.Type.Id != currencyCodeTo {
		var requestError *http.Error
		requestError, currencyTo, err = exchange.Convert(currencyFrom, currencyCodeTo, cfg.GetString("secretKey"))
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR)
		}
		if requestError != nil {
			log.ErrorLogger.Println(requestError.Serialize())
			return http.ToError(http.INTERNAL_SERVER_ERROR)
		}
	} else {
		currencyTo = currencyFrom
	}

	// Check balance
	balance, err := accounts.GetBalance(senderId, currencyFrom.Type.Id)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if balance.Cmp(currencyFrom.Amount) < 0 {
		return &http.Error{
			Code:    http.INSUFFICIENT_FUNDS,
			Message: "Insufficient Funds, Top Up Your Balance First",
		}
	}

	// Process payment
	payment := payments.New(senderId, receiverId, currencyFrom, currencyTo, timestamp)

	err = payment.Commit()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return nil
}

func addPayment(receiverId uint64, currency *currency.Currency, timestamp uint64) *http.Error {

	exists, err := accounts.Exists(receiverId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !exists {
		return &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Receiver Not Found!",
		}
	}

	payment := payments.New(0, receiverId, currency, currency, timestamp)

	err = payment.Commit()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return nil
}

func getBalance(accountId uint64, currencyCode uint64) (*http.Error, *big.Int) {
	if currency.FromCode(currencyCode) == nil {
		return http.ToError(http.INVALID_CURRENCY_CODE), nil
	}
	exists, err := accounts.Exists(accountId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	if !exists {
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

func getHistory(accountId uint64) (*http.Error, *[]payments.Payment) {
	/* TODO
	exists, err := accounts.Exists(accountId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	if !exists {
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
	return nil, &history*/return nil, nil
}

func getAccountInfo(accountId uint64) (*http.Error, *accounts.AccountInfo) {
	accountInfo, err := accounts.GetAccountInfo(accountId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR), nil
	}
	return nil, accountInfo
}

func getPayment(id uint64, accountId uint64) (*http.Error, *payments.Payment) {
	payment := payments.Payment{
		Id: id,
	}
	var currencyCode uint64
	var amountBase64 string
	err, found := database.Payments.Get( // TODO
		accountId,
		eplidr.Keys{{"id", id}},
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

func createAccount(username string, timestamp uint64) *http.Error {
	accountId := fnv64(username)
	_, found, err := database.Accounts.GetString(accountId, "name")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if found {
		return &http.Error{
			Code:    http.ALREADY_EXISTS,
			Message: "Account For User Already Created!",
		}
	}
	err = database.Accounts.Put(accountId, []string{"name", "id"}, []interface{}{username, accountId})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	go func() {
		for _, Type := range currency.Types {
			err = database.Balances.Put(accountId, eplidr.Columns{{"id", accountId}, {"balance", utils.EncodeBigInt(new(big.Int).SetInt64(0))}, {"onSerial", timestamp}, {"currency", Type.Id}})
			if err != nil {
				log.ErrorLogger.Println(err.Error())
			}
		}
	}()
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
