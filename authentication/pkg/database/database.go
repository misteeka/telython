package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"telython/authentication/pkg/cfg"
	"telython/authentication/pkg/log"
	"telython/pkg/eplidr"
	"time"
)

var db *sql.DB

var (
	Users       *eplidr.Table
	SingleUsers *eplidr.SingleKeyTable

	EmailCodes                *eplidr.SingleKeyTable
	PendingEmailConfirmations *eplidr.SingleKeyTable
)

func InitDatabase() error {
	var err error
	dataSource := "{user}:{password}@tcp(localhost:41091)/{db}"
	dataSource = strings.Replace(dataSource, "{user}", cfg.GetString("user"), 1)
	dataSource = strings.Replace(dataSource, "{password}", cfg.GetString("password"), 1)
	dataSource = strings.Replace(dataSource, "{db}", cfg.GetString("dbName"), 1)
	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetMaxIdleConns(cfg.GetInt("maxIdleConns"))
	db.SetMaxOpenConns(cfg.GetInt("maxOpenConns"))

	Users, err = eplidr.NewTable(
		"users",
		1,
		[]string{
			`CREATE TABLE IF NOT EXISTS {table} (
				name varchar(255) primary key {nn},
				password varchar(44) {nn},
				email varchar(255) {nn},
				reg_date uint64 {nn},
				last_login uint64 {nn},
				last_ip varchar(40) {nn},
				reg_ip varchar(40) {nn}
			);`,
			"CREATE INDEX index_email ON users(email);",
		},
		db,
	)
	if err != nil {
		return err
	}
	SingleUsers = eplidr.SingleKeyImplementation(Users, "name")
	EmailCodes, err = eplidr.NewSingleKeyTable(
		"emailcodes",
		"name",
		1,
		[]string{
			`CREATE TABLE IF NOT EXISTS {table} (
				name varchar(255) {nn} primary key,
				code int {nn}
			);`,
		},
		db,
	)
	if err != nil {
		return err
	}
	PendingEmailConfirmations, err = eplidr.NewSingleKeyTable(
		"pending_email_confirmations",
		"name",
		1,
		[]string{
			`CREATE TABLE IF NOT EXISTS {table} (
				name varchar(255) {nn} primary key,
				email varchar(255) {nn},
				code int {nn},
				timestamp uint64 {nn}
			);`,
			"CREATE INDEX index_email ON pending_email_confirmations (email);",
		},
		db,
	)
	if err != nil {
		return err
	}
	return nil
}

func Exec(sql string, v ...any) (sql.Result, error) {
	return db.Exec(sql, v)
}

func Query(sql string, v ...any) (*sql.Rows, error) {
	return db.Query(sql, v)
}

func ReleaseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.ErrorLogger.Print(err.Error())
		return
	}
}