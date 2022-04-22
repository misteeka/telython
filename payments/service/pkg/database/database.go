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
	Accounts   *eplidr.SingleKeyTable
	Balances   *eplidr.Table
	Payments   *eplidr.Table
	LastActive *eplidr.SingleKeyTable
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

	Payments, err = eplidr.NewTable(
		"payments",
		4,
		[]string{
			"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn},`sender` uint64 {nn},`receiver` uint64 {nn},`amountFrom` varchar(32) {nn},`amountTo` varchar(32) {nn},`currencyFrom` uint64 {nn},`currencyTo` uint64 {nn},`timestamp` uint64 {nn});",
			"create index index_sender on {table} (sender);",
			"create index index_receiver on {table} (receiver);",
			"create index index_serial on {table} (timestamp);",
			"create index index_currency on {table} (currencyTo);",
		},
		defaultDriver,
	)
	if err != nil {
		return err
	}
	Accounts, err = eplidr.NewSingleKeyTable(
		"accounts",
		"id",
		4,
		[]string{"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 primary key {nn}, `name` varchar(255) {nn});"},
		defaultDriver,
	)
	if err != nil {
		return err
	}
	Balances, err = eplidr.NewTable(
		"balances",
		4,
		[]string{
			"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn}, `balance` varchar(20) {nn}, `onSerial` uint64 {nn}, `currency` uint64 {nn});",
			"create index index_id on {table} (id);",
			"create index index_onSerial on {onSerial} (id);",
			"create index index_currency on {currency} (id);",
		},
		defaultDriver,
	)
	if err != nil {
		return err
	}

	//Payments.DropUnsafe()
	//Accounts.Table.DropUnsafe()
	//Balances.Table.DropUnsafe()
	return nil
}

func ReleaseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.ErrorLogger.Print(err.Error())
		return
	}
}
