package accounts

import (
	"fmt"
	"telython/payments/service/pkg/database"
)

type Account struct {
	Id       uint64
	Currency int
}

func Load(id uint64) (*Account, error) {
	rows, err := database.Accounts.Query(fmt.Sprintf("SELECT `currency` FROM {table} WHERE `id` = %d;", id), id)
	if err != nil {
		return nil, err
	}
	var account Account
	account.Id = id
	if rows.Next() {
		err = rows.Scan(&account.Currency)
		if err != nil {
			rows.Close()
			return nil, err
		}
	} else {
		return nil, nil
	}
	rows.Close()
	return &account, nil
}

// SELECT `amount` FROM `payments2` WHERE `sender` = 15381326603262689376 AND `serial` > 0

func (account *Account) GetBalance() (uint64, error) {
	rows, err := database.Balances.Query(fmt.Sprintf("SELECT `balance`, `onSerial` FROM {table} WHERE `id` = %d;", account.Id), account.Id)
	if err != nil {
		return 0, err
	}
	var balance uint64
	var onSerial uint64
	if rows.Next() {
		err = rows.Scan(&balance, &onSerial)
		if err != nil {
			rows.Close()
			return 0, err
		}
	} else {
		return 0, nil
	}
	rows.Close()
	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amount` FROM {table} WHERE `sender` = %d AND `timestamp` > %d", account.Id, onSerial), account.Id)
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
	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amount` FROM {table} WHERE `receiver` = %d AND `timestamp` > %d", account.Id, onSerial), account.Id)
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

func (account *Account) FullScanBalance() (uint64, error) {
	var balance uint64
	rows, err := database.Payments.Query(fmt.Sprintf("SELECT `amount` FROM {table} WHERE `sender` = %d", account.Id), account.Id)
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
	rows, err = database.Payments.Query(fmt.Sprintf("SELECT `amount` FROM {table} WHERE `receiver` = %d", account.Id), account.Id)
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

func (account *Account) RecalculateBalance(timestamp uint64) (uint64, error) {
	balance, err := account.FullScanBalance()
	if err != nil {
		return 0, err
	}
	database.Balances.Put(account.Id, []string{"onSerial", "balance"}, []interface{}{timestamp, balance})
	return balance, nil
}
