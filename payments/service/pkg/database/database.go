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
	Balances   *eplidr.SingleKeyTable
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
			"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn},`sender` uint64 {nn},`receiver` uint64 {nn},`amount` varchar(16)  {nn},`currency` uint64 {nn},`timestamp` uint64 {nn});",
			"create index index_sender on {table} (sender);",
			"create index index_receiver on {table} (receiver);",
			"create index index_serial on {table} (timestamp);",
		},
		defaultDriver,
	)
	if err != nil {
		return err
	}
	Accounts, err = eplidr.NewSingleKeyTable(
		"accounts",
		"name",
		4,
		[]string{"CREATE TABLE IF NOT EXISTS {table} (`nameHash` uint64 primary key {nn}, `name` varchar(255) {nn});"},
		defaultDriver,
	)
	if err != nil {
		return err
	}
	Balances, err = eplidr.NewSingleKeyTable(
		"balances",
		"id",
		4,
		[]string{"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn} primary key, `balance` varchar(16) {nn} default 0, `onSerial` uint64 {nn} default 0);"},
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
