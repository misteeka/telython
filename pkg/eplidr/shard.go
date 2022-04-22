package eplidr

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Shard struct {
	table  *Table
	driver *sql.DB
	num    uint
}

func (shard *Shard) GetString(key Key, column string) (string, bool, error) {
	var result string
	err, found := shard.Get(Keys{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return "", found, err
	}
	return result, found, nil
}
func (shard *Shard) GetInt(key Key, column string) (int, bool, error) {
	var result int
	err, found := shard.Get(Keys{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetInt64(key Key, column string) (int64, bool, error) {
	var result int64
	err, found := shard.Get(Keys{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetFloat(key Key, column string) (float64, bool, error) {
	var result float64
	err, found := shard.Get(Keys{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetUint64(key Key, column string) (uint64, bool, error) {
	var result uint64
	err, found := shard.Get(Keys{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetUint(key Key, column string) (uint, bool, error) {
	var result uint
	err, found := shard.Get(Keys{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetBoolean(key Key, column string) (bool, bool, error) {
	var result bool
	err, found := shard.Get(Keys{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return false, found, err
	}
	return result, found, nil
}

func (shard *Shard) Get(keys Keys, columnNames []string, data []interface{}) (error, bool) {
	query := fmt.Sprintf("SELECT %s FROM {table} %s;", ColumnNamesToQuery(columnNames...), keys.Query())
	rows, err := shard.Query(query)
	if err != nil {
		return err, false
	}
	if rows.Next() {
		err := rows.Scan(data...)
		if err != nil {
			rows.Close()
			return err, true
		}
		rows.Close()
	} else {
		rows.Close()
		return nil, false
	}
	return nil, true
}
func (shard *Shard) Put(values Columns) error {
	// `%s` = ?
	columnsString := ""
	valuesString := ""
	for i := 0; i < len(values); i++ {
		if i == len(values)-1 {
			columnsString += fmt.Sprintf("`%s`", values[i].Name)
		} else {
			columnsString += fmt.Sprintf("`%s`, ", values[i].Name)
		}
	}
	for i := 0; i < len(values); i++ {
		if i == len(values)-1 {
			valuesString += fmt.Sprintf("%s", value(values[i].Value))
		} else {
			valuesString += fmt.Sprintf("%s, ", value(values[i].Value))
		}
	}
	query := fmt.Sprintf("INSERT INTO {table} (%s) values (%s);", columnsString, valuesString)
	_, err := shard.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (shard *Shard) Set(keys Keys, values Columns) error {
	s := ""
	for i := 0; i < len(values); i++ {
		if i == len(values)-1 {
			s += fmt.Sprintf("`%s` = %s", values[i].Name, value(values[i].Value))
		} else {
			s += fmt.Sprintf("`%s` = %s, ", values[i].Name, value(values[i].Value))
		}
	}
	query := fmt.Sprintf("UPDATE {table} SET %s %s;", s, keys.Query())
	_, err := shard.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (shard *Shard) Add(keys Keys, values Columns) error {
	s := ""
	for i := 0; i < len(values); i++ {
		if i == len(values)-1 {
			s += fmt.Sprintf("`%s` = `%s` + %s", values[i].Name, values[i].Name, value(values[i].Value))
		} else {
			s += fmt.Sprintf("`%s` = `%s` + %s, ", values[i].Name, values[i].Name, value(values[i].Value))
		}
	}
	query := fmt.Sprintf("UPDATE {table} SET %s %s;", s, keys.Query())
	_, err := shard.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (shard *Shard) Remove(keys Keys) error {
	query := fmt.Sprintf("DELETE FROM {table} %s;", keys.Query())
	_, err := shard.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (shard *Shard) Exec(query string) (sql.Result, error) {
	query = strings.Replace(query, "{table}", fmt.Sprintf("`%s`", shard.table.GetName(shard.num)), 1)
	return shard.driver.Exec(query)
}
func (shard *Shard) Query(query string) (*sql.Rows, error) {
	query = strings.Replace(query, "{table}", fmt.Sprintf("`%s`", shard.table.GetName(shard.num)), 1)
	return shard.driver.Query(query)
}

func (shard *Shard) ReleaseRows(rows *sql.Rows) error {
	return rows.Close()
}

func (shard *Shard) RawTx() (*sql.Tx, error) {
	return shard.driver.Begin()
}

func (shard *Shard) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return shard.driver.BeginTx(ctx, opts)
}

func (shard *Shard) SingleSet(keys Keys, column Column) error {
	return shard.Set(keys, Columns{column})
}

func (shard *Shard) Drop() error {
	_, err := shard.driver.Exec(fmt.Sprintf("DROP TABLE %s;", shard.table.GetName(shard.num)))
	return err
}
