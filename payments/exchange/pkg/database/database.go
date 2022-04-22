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
	Prices *eplidr.SingleKeyTable
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

	Prices, err = eplidr.NewSingleKeyTable(
		"prices",
		"symbol",
		2,
		[]string{"CREATE TABLE IF NOT EXISTS {table} (`symbol` varchar(5) primary key, `price` varchar(32))"},
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
