package accounts

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"sort"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/database"
	"telython/pkg/eplidr"
	"telython/pkg/log"
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

func GetBalance(username string, currencyCode uint64) (*big.Int, error) {
	accountId := fnv64(username)
	rows, err := database.Balances.Query(fmt.Sprintf("SELECT `balance`, `onSerial` FROM {table} WHERE `id` = %d AND `currency` = %d;", fnv64(username), currencyCode), accountId)
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
		var ok bool
		balance, ok = new(big.Int).SetString(balanceString, 10)
		if !ok {
			balance = big.NewInt(0)
		}
	} else {
		changed = true
		notFound = true
		balance = big.NewInt(0)
	}
	rows.Close()
	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amountTo` FROM {table} WHERE `receiver` = '%s' AND `timestamp` > %d AND `currencyTo` = %d;", username, timestamp, currencyCode), accountId)
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
		amount, ok := new(big.Int).SetString(amountString, 10)
		if !ok {
			continue
		}
		balance.Add(balance, amount)
	}
	rows.Close()

	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amountFrom` FROM {table} WHERE `sender` = '%s' AND `timestamp` > %d AND `currencyFrom` = %d;", username, timestamp, currencyCode), accountId)
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
		amount, ok := new(big.Int).SetString(amountString, 10)
		if !ok {
			continue
		}
		balance.Sub(balance, amount)
	}
	rows.Close()
	if changed {
		go func() {
			if notFound {
				err := database.Balances.Put(accountId, eplidr.Columns{{"id", accountId}, {"balance", balance.String()}, {"onSerial", time.Now().UnixMicro()}, {"currency", currencyCode}})
				if err != nil {
					log.ErrorLogger.Println(err.Error())
					return
				}
			} else {
				err := database.Balances.Set(accountId, eplidr.Keys{{"id", accountId}, {"currency", currencyCode}}, eplidr.Columns{{"balance", balance.String()}, {"onSerial", time.Now().UnixMicro()}})
				if err != nil {
					log.ErrorLogger.Println(err.Error())
					return
				}
			}
		}()
	}
	return balance, nil
}

// FullScanBalance() (*big.Int, error) {

func Exists(username string) (bool, error) {
	return ExistsId(fnv64(username))
}

func ExistsId(accountId uint64) (bool, error) {
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
	Id       uint64
	balances []currency.Currency
}

func (accountInfo *AccountInfo) Serialize() string {
	serializedBalances := ""
	for i := 0; i < len(accountInfo.balances); i++ {
		if i == len(accountInfo.balances)-1 {
			serializedBalances += `"` + base64.StdEncoding.EncodeToString([]byte(accountInfo.balances[i].Json())) + `"`
		} else {
			serializedBalances += `"` + base64.StdEncoding.EncodeToString([]byte(accountInfo.balances[i].Json())) + `",`
		}
	}
	return fmt.Sprintf(`{"id": %d, "balances": [%s]}`, accountInfo.Id, serializedBalances)
}

func GetAccountInfo(username string) (*AccountInfo, error) {
	accountId := fnv64(username)
	var balances []currency.Currency
	for _, v := range currency.Types {
		amount, err := GetBalance(username, v.Id)
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
	sort.Slice(balances[:], func(i, j int) bool {
		return balances[i].Type.Id < balances[j].Type.Id
	})
	return &AccountInfo{
		Id:       accountId,
		balances: balances,
	}, nil
}
