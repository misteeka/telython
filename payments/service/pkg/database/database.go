package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"telython/payments/service/pkg/cfg"
	"telython/payments/service/pkg/log"
	"telython/pkg/eplidr"
	"time"
)

var (
	Accounts   *eplidr.Table
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
			"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn},`sender` uint64 {nn},`receiver` uint64 {nn},`amount` uint64 {nn},`currency` int {nn},`timestamp` uint64 {nn});",
			"create index index_sender on {table} (sender);",
			"create index index_receiver on {table} (receiver);",
			"create index index_serial on {table} (timestamp);",
		},
		defaultDriver,
	)
	if err != nil {
		return err
	}
	Accounts, err = eplidr.NewTable(
		"accounts",
		4,
		[]string{"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn} primary key, `name` varchar(255) {nn}, `currency` int default 0 {nn});"},
		defaultDriver,
	)
	if err != nil {
		return err
	}
	Balances, err = eplidr.NewSingleKeyTable(
		"balances",
		"id",
		4,
		[]string{"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn} primary key, `balance` uint64 {nn} default 0, `onSerial` uint64 {nn} default 0);"},
		defaultDriver,
	)
	if err != nil {
		return err
	}

	/*LastActive = eplidr.NewSingleKeyTable(
		"lastactive",
		"name",
		1,
		"",
		defaultDriver,
	)*/
	return nil
}

func ReleaseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.ErrorLogger.Print(err.Error())
		return
	}
}
