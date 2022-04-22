package accounts

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/database"
	"telython/pkg/eplidr"
	"telython/pkg/log"
	"telython/pkg/utils"
	"time"
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

func GetBalance(accountId uint64, currencyCode uint64) (*big.Int, error) {
	rows, err := database.Balances.Query(fmt.Sprintf("SELECT `balance`, `onSerial` FROM {table} WHERE `id` = %d AND `currency` = %d;", accountId, currencyCode), accountId)
	if err != nil {
		return nil, err
	}
	var balanceString string
	var balance *big.Int
	var timestamp uint64
	var changed bool
	var notFound bool
	if rows.Next() {
		err = rows.Scan(&balanceString, &timestamp)
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
		changed = true
		notFound = true
		balance = big.NewInt(0)
	}
	rows.Close()
	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amountTo` FROM {table} WHERE `receiver` = %d AND `timestamp` > %d AND `currencyTo` = %d;", accountId, timestamp, currencyCode), accountId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		changed = true
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

	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amountFrom` FROM {table} WHERE `sender` = %d AND `timestamp` > %d AND `currencyFrom` = %d;", accountId, timestamp, currencyCode), accountId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		changed = true
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
	if changed {
		go func() {
			if notFound {
				err := database.Balances.Put(accountId, eplidr.Columns{{"id", accountId}, {"balance", utils.EncodeBigInt(balance)}, {"onSerial", time.Now().UnixMicro()}, {"currency", currencyCode}})
				if err != nil {
					log.ErrorLogger.Println(err.Error())
					return
				}
			} else {
				err := database.Balances.Set(accountId, eplidr.Keys{{"id", accountId}, {"currency", currencyCode}}, eplidr.Columns{{"balance", utils.EncodeBigInt(balance)}, {"onSerial", time.Now().UnixMicro()}})
				if err != nil {
					log.ErrorLogger.Println(err.Error())
					return
				}
			}
		}()
	}
	return balance, nil
}

func FullScanBalance(accountId uint64, currencyCode uint64) (uint64, error) {
	var balance uint64
	rows, err := database.Payments.Query(fmt.Sprintf("SELECT `amountTo` FROM {table} WHERE `receiver` = %d AND `currencyTo` = %d;", accountId, currencyCode), accountId)
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
	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amountFrom` FROM {table} WHERE `sender` = %d AND `currencyFrom` = %d;", accountId, currencyCode), accountId)
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
	return balance, nil
}

func Exists(accountId uint64) (bool, error) {
	_, found, err := database.Accounts.GetString(accountId, "name")
	if err != nil {
		return false, err
	}
	return found, nil
}

func GetId(username string) (uint64, bool, error) {
	id := fnv64(username)
	_, exists, err := database.Accounts.GetString(id, "name")
	return id, exists, err
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
