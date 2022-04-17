package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"telython/payments/gates/ethgate/pkg/cfg"
	"telython/payments/gates/ethgate/pkg/log"
	"telython/pkg/eplidr"
	"time"
)

var (
	Accounts *eplidr.SingleKeyTable
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

	Accounts, err = eplidr.NewSingleKeyTable(
		"accounts",
		"id",
		4,
		[]string{"CREATE TABLE IF NOT EXISTS {table} (`id` uint64 {nn} primary key, `public` varchar(128) {nn}, `private` varchar(128) {nn});"},
		defaultDriver,
	)
	if err != nil {
		return err
	}
	return nil
}

func ReleaseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.ErrorLogger.Print(err.Error())
		return
	}
}
