package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"telython/pkg/cfg"
	"telython/pkg/eplidr"
	"telython/pkg/log"
	"time"
)

var (
	AccountToWallet *eplidr.SingleKeyTable
	WalletToAccount *eplidr.SingleKeyTable
)

func InitDatabase() error {
	var err error
	dataSource := "{user}:{password}@tcp(localhost:41091)/{db}"
	dataSource = strings.Replace(dataSource, "{user}", cfg.GetString("user"), 1)
	dataSource = strings.Replace(dataSource, "{password}", cfg.GetString("password"), 1)
	dataSource = strings.Replace(dataSource, "{db}", cfg.GetString("dbName"), 1)
	defaultDriver, err := sql.Open("mysql", dataSource)
	if err != nil {
		return err
	}
	defaultDriver.SetConnMaxLifetime(0)
	defaultDriver.SetConnMaxIdleTime(1 * time.Minute)
	defaultDriver.SetMaxIdleConns(cfg.GetInt("maxIdleConns"))
	defaultDriver.SetMaxOpenConns(cfg.GetInt("maxOpenConns"))

	AccountToWallet, err = eplidr.NewSingleKeyTable(
		"account_to_wallet",
		"id",
		2,
		[]string{
			"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn} primary key, `address` varchar(128) {nn}, `private` varchar(128) {nn});",
		},
		defaultDriver,
	)
	if err != nil {
		return err
	}
	WalletToAccount, err = eplidr.NewSingleKeyTable(
		"wallet_to_account",
		"address",
		2,
		[]string{
			"CREATE TABLE IF NOT EXISTS {table} (`address` varchar(128) {nn} primary key, `name` varchar(128) {nn});",
		},
		defaultDriver,
	)

	//WalletToAccount.Table.DropUnsafe()
	//AccountToWallet.Table.DropUnsafe()
	return nil
}

func ReleaseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.ErrorLogger.Print(err.Error())
		return
	}
}
