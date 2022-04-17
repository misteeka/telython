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

func (shard *Shard) GetString(key Column, column string) (string, bool, error) {
	var result string
	err, found := shard.Get(Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return "", found, err
	}
	return result, found, nil
}
func (shard *Shard) GetInt(key Column, column string) (int, bool, error) {
	var result int
	err, found := shard.Get(Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetInt64(key Column, column string) (int64, bool, error) {
	var result int64
	err, found := shard.Get(Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetFloat(key Column, column string) (float64, bool, error) {
	var result float64
	err, found := shard.Get(Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetUint64(key Column, column string) (uint64, bool, error) {
	var result uint64
	err, found := shard.Get(Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetUint(key Column, column string) (uint, bool, error) {
	var result uint
	err, found := shard.Get(Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (shard *Shard) GetBoolean(key Column, column string) (bool, bool, error) {
	var result bool
	err, found := shard.Get(Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return false, found, err
	}
	return result, found, nil
}

func (shard *Shard) Get(keys Columns, columnNames []string, data []interface{}) (error, bool) {
	query := fmt.Sprintf("SELECT %s FROM {table} %s;", ColumnNamesToQuery(columnNames...), KeysToQuery(keys))
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
			columnsString += fmt.Sprintf("`%s`", values[i].Key)
		} else {
			columnsString += fmt.Sprintf("`%s`, ", values[i].Key)
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
func (shard *Shard) Set(keys Columns, values Columns) error {
	s := ""
	for i := 0; i < len(values); i++ {
		if i == len(values)-1 {
			s += fmt.Sprintf("`%s` = %s", values[i].Key, value(values[i].Value))
		} else {
			s += fmt.Sprintf("`%s` = %s, ", values[i].Key, value(values[i].Value))
		}
	}
	query := fmt.Sprintf("UPDATE {table} SET %s %s;", s, KeysToQuery(keys))
	_, err := shard.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (shard *Shard) Remove(keys Columns) error {
	query := fmt.Sprintf("DELETE FROM `{table}` %s;", KeysToQuery(keys))
	_, err := shard.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (shard *Shard) Exec(query string) (sql.Result, error) {
	//log.InfoLogger.Println(query)
	query = strings.Replace(query, "{table}", fmt.Sprintf("`%s`", shard.table.GetName(shard.num)), 1)
	return shard.driver.Exec(query)
}
func (shard *Shard) Query(query string) (*sql.Rows, error) {
	//log.InfoLogger.Println(query)
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

func (shard *Shard) SingleSet(keys Columns, column Column) error {
	return shard.Set(keys, Columns{column})
}

func (shard *Shard) Drop() error {
	_, err := shard.driver.Exec(fmt.Sprintf("DROP TABLE %s;", shard.table.GetName(shard.num)))
	return err
}
