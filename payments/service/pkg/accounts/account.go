package accounts

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/database"
)

// SELECT `amount` FROM `payments2` WHERE `sender` = 15381326603262689376 AND `serial` > 0

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

func GetBalance(nameHash uint64, currencyCode uint64) (*big.Int, error) {
	rows, err := database.Balances.Query(fmt.Sprintf("SELECT `balance`, `onSerial` FROM {table} WHERE `id` = %d AND `currency` = %d;", nameHash, currencyCode), nameHash)
	if err != nil {
		return nil, err
	}
	var balanceString string
	var balance *big.Int
	var onSerial uint64
	if rows.Next() {
		err = rows.Scan(&balanceString, &onSerial)
		if err != nil {
			rows.Close()
			return nil, err
		}
		balanceBytes, err := base64.StdEncoding.DecodeString(balanceString)
		if err != nil {
			rows.Close()
			return nil, err
		}
		balance = new(big.Int).SetBytes(balanceBytes)
	} else {
		balance = big.NewInt(0)
	}
	rows.Close()
	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amount` FROM {table} WHERE `receiver` = %d AND `timestamp` > %d", nameHash, onSerial), nameHash)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var amountString string
		err = rows.Scan(&amountString)
		if err != nil {
			rows.Close()
			return nil, err
		}
		amountBytes, err := base64.StdEncoding.DecodeString(amountString)
		if err != nil {
			rows.Close()
			return nil, err
		}
		balance.Add(balance, new(big.Int).SetBytes(amountBytes))
	}
	rows.Close()

	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amount` FROM {table} WHERE `sender` = %d AND `timestamp` > %d", nameHash, onSerial), nameHash)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var amountString string
		err = rows.Scan(&amountString)
		if err != nil {
			rows.Close()
			return nil, err
		}
		amountBytes, err := base64.StdEncoding.DecodeString(amountString)
		if err != nil {
			rows.Close()
			return nil, err
		}
		balance.Sub(balance, new(big.Int).SetBytes(amountBytes))
	}
	rows.Close()
	return balance, nil
}

func FullScanBalance(nameHash uint64, currencyCode uint64) (uint64, error) {
	var balance uint64
	rows, err := database.Payments.Query(fmt.Sprintf("SELECT `amount` FROM {table} WHERE `sender` = %d AND `currency` = %d;", nameHash, currencyCode), nameHash)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var amount uint64
		err = rows.Scan(&amount)
		if err != nil {
			rows.Close()
			return 0, err
		}
		balance -= amount
	}
	rows.Close()
	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amount` FROM {table} WHERE `receiver` = %d AND `currency` = %d;", nameHash, currencyCode), nameHash)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var amount uint64
		err = rows.Scan(&amount)
		if err != nil {
			rows.Close()
			return 0, err
		}
		balance += amount
	}
	rows.Close()
	return balance, nil
}

func GetId(username string) (uint64, bool, error) {
	return database.Accounts.GetUint(fnv64(username), "nameHash")
}

type AccountInfo struct {
	nameHash uint64
	balances *[]currency.Currency
}

func GetAccountInfo(nameHash uint64) (*AccountInfo, error) {
	var balances []currency.Currency
	for _, v := range currency.Types {
		amount, err := GetBalance(nameHash, v.Id)
		if err != nil {
			balances = append(balances, currency.Currency{
				Type:   v,
				Amount: nil,
			})
			continue
		}
		balances = append(balances, currency.Currency{
			Type:   v,
			Amount: amount,
		})
	}
	return &AccountInfo{
		nameHash: nameHash,
		balances: &balances,
	}, nil
}
