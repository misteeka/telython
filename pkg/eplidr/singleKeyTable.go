package eplidr

import (
	"context"
	"database/sql"
)

type SingleKeyTable struct {
	Table *Table
	key   string
}

func NewSingleKeyTable(name string, key string, shardsCount uint, creatingQuery []string, drivers Drivers) (*SingleKeyTable, error) {
	// params:
	// [0] dataSource
	// [1]
	table, err := NewTable(name, shardsCount, creatingQuery, drivers)
	if err != nil {
		return nil, err
	}
	return &SingleKeyTable{
		Table: table,
		key:   key,
	}, nil
}

func SingleKeyImplementation(keyTable *Table, key string) *SingleKeyTable {
	return &SingleKeyTable{
		Table: keyTable,
		key:   key,
	}
}

func (table *SingleKeyTable) GetString(key interface{}, column string) (string, bool, error) {
	return table.Table.GetString(Column{Key: table.key, Value: key}, column)
}
func (table *SingleKeyTable) GetInt(key interface{}, column string) (int, bool, error) {
	return table.Table.GetInt(Column{Key: table.key, Value: key}, column)
}
func (table *SingleKeyTable) GetInt64(key interface{}, column string) (int64, bool, error) {
	return table.Table.GetInt64(Column{Key: table.key, Value: key}, column)
}
func (table *SingleKeyTable) GetFloat(key interface{}, column string) (float64, bool, error) {
	return table.Table.GetFloat(Column{Key: table.key, Value: key}, column)
}
func (table *SingleKeyTable) GetUint(key interface{}, column string) (uint64, bool, error) {
	return table.Table.GetUint(Column{Key: table.key, Value: key}, column)
}
func (table *SingleKeyTable) GetBoolean(key interface{}, column string) (bool, bool, error) {
	return table.Table.GetBoolean(Column{Key: table.key, Value: key}, column)
}

func (table *SingleKeyTable) Get(key interface{}, columns []string, data []interface{}) (error, bool) {
	return table.Table.Get(key, Columns{Column{Key: table.key, Value: key}}, columns, data)
}
func (table *SingleKeyTable) Set(key interface{}, columns []string, values []interface{}) error {
	return table.Table.Set(key, Columns{Column{Key: table.key, Value: key}}, PlainToColumns(columns, values))
}

func (table *SingleKeyTable) SingleSet(key interface{}, column string, value interface{}) error {
	return table.Table.Set(key, Columns{Column{Key: table.key, Value: key}}, Columns{Column{
		Key:   column,
		Value: value,
	}})
}

func (table *SingleKeyTable) Put(key interface{}, columns []string, values []interface{}) error {
	return table.Table.Put(key, PlainToColumns(columns, values))
}

func (table *SingleKeyTable) Remove(key interface{}) error {
	return table.Table.Remove(table.key, Columns{Column{
		Key:   table.key,
		Value: key,
	}})
}

func (table *SingleKeyTable) ReleaseRows(rows *sql.Rows) error {
	return rows.Close()
}

func (table *SingleKeyTable) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return nil, nil
}

func (table *SingleKeyTable) Exec(query string, key interface{}) (sql.Result, error) {
	return table.Table.Exec(query, key)
}
func (table *SingleKeyTable) Query(query string, key interface{}) (*sql.Rows, error) {
	return table.Table.Query(query, key)
}
