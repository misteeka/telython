package eplidr

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Table struct {
	name        string
	shardsCount uint
	Shards      []*Shard

	creatingQuery []string

	hashFunc func(interface{}) uint
}

type Drivers interface{}

func NewTable(name string, shardsCount uint, creatingQuery []string, driverParam Drivers) (*Table, error) {
	var table *Table
	switch dataSource := driverParam.(type) {
	case []*sql.DB:
		table = &Table{
			name:          name,
			shardsCount:   shardsCount,
			creatingQuery: creatingQuery,
			hashFunc:      StandardGetShardFunc,
		}
		shards := make([]*Shard, len(dataSource))
		for i := 0; i < len(dataSource); i++ {
			shards[i] = &Shard{
				table:  table,
				driver: dataSource[i],
				num:    uint(i),
			}
		}
		table.Shards = shards
	case *sql.DB:
		drivers := make([]*sql.DB, shardsCount)
		for i := 0; i < int(shardsCount); i++ {
			drivers[i] = dataSource
		}
		table = &Table{
			name:          name,
			shardsCount:   shardsCount,
			creatingQuery: creatingQuery,
			hashFunc:      StandardGetShardFunc,
		}
		shards := make([]*Shard, len(drivers))
		for i := 0; i < len(drivers); i++ {
			shards[i] = &Shard{
				table:  table,
				driver: drivers[i],
				num:    uint(i),
			}
		}
		table.Shards = shards
	}
	return table, table.init()
}

func Serialize(value interface{}) string {
	switch array := value.(type) {
	case map[string]interface{}:
		var result string
		keys := make([]string, 0, len(array))
		values := make([]interface{}, 0, len(array))
		for i := 0; i < len(array); i++ {
			if i == len(array)-1 {
				result += fmt.Sprintf("%s:%v", keys[i], values[i])
			} else {
				result += fmt.Sprintf("%s:%v,", keys[i], values[i])
			}
		}
		return result
	case []interface{}:
		var result string
		for i := 0; i < len(array); i++ {
			if i == len(array)-1 {
				result += fmt.Sprintf("%v", array[i])
			} else {
				result += fmt.Sprintf("%v,", array[i])
			}
		}
		return result
	default:
		return ""
	}
}
func DeserializeMap(serializedData string) {
	var result map[string]interface{}
	for i := 0; i < len(serializedData); i++ {
		key := ""
		value := ""
		result[key] = value
	}
}
func DeserializeSlice(serializedData string) {
	for i := 0; i < len(serializedData); i++ {

	}

}

func (table *Table) GetName(shard uint) string {
	return table.name + strconv.FormatUint(uint64(shard), 10)
}
func (table *Table) GetShardNum(key interface{}) uint {
	return table.hashFunc(key) % table.shardsCount
}
func (table *Table) GetShard(num uint) *Shard {
	return table.Shards[num]
}

func (table *Table) init() error {
	for i := 0; i < len(table.creatingQuery); i++ {
		table.creatingQuery[i] = strings.ReplaceAll(table.creatingQuery[i], "uint64", "BIGINT UNSIGNED")
		table.creatingQuery[i] = strings.ReplaceAll(table.creatingQuery[i], "int", "INTEGER")
		table.creatingQuery[i] = strings.ReplaceAll(table.creatingQuery[i], "{nn}", "NOT NULL")
		table.creatingQuery[i] = strings.ReplaceAll(table.creatingQuery[i], "{n}", "NULL")
		for a := 0; a < len(table.Shards); a++ {
			rows, err := table.Shards[a].Query(fmt.Sprintf("SHOW TABLES LIKE '%s';", table.GetName(uint(a))))
			if err != nil {
				return err
			}
			if rows.Next() {
				rows.Close()
				continue
			}
			_, err = table.Shards[i].Exec(strings.Replace(table.creatingQuery[i], "{table}", fmt.Sprintf("`%s`", table.GetName(uint(a))), 1))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (table *Table) GetStringMultiple(shardKey interface{}, keys Columns, column string) (string, bool, error) {
	var result string
	err, found := table.Get(shardKey, keys, []string{column}, []interface{}{&result})
	if err != nil {
		return "", found, err
	}
	return result, found, nil
}
func (table *Table) GetIntMultiple(shardKey interface{}, keys Columns, column string) (int, bool, error) {
	var result int
	err, found := table.Get(shardKey, keys, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *Table) GetInt64Multiple(shardKey interface{}, keys Columns, column string) (int64, bool, error) {
	var result int64
	err, found := table.Get(shardKey, keys, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *Table) GetFloatMultiple(shardKey interface{}, keys Columns, column string) (float64, bool, error) {
	var result float64
	err, found := table.Get(shardKey, keys, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *Table) GetUintMultiple(shardKey interface{}, keys Columns, column string) (uint64, bool, error) {
	var result uint64
	err, found := table.Get(shardKey, keys, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *Table) GetBooleanMultiple(shardKey interface{}, keys Columns, column string) (bool, bool, error) {
	var result bool
	err, found := table.Get(shardKey, keys, []string{column}, []interface{}{&result})
	if err != nil {
		return false, found, err
	}
	return result, found, nil
}

func (table *Table) GetString(key Column, column string) (string, bool, error) {
	var result string
	err, found := table.Get(key.Value, Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return "", found, err
	}
	return result, found, nil
}
func (table *Table) GetInt(key Column, column string) (int, bool, error) {
	var result int
	err, found := table.Get(key.Value, Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *Table) GetInt64(key Column, column string) (int64, bool, error) {
	var result int64
	err, found := table.Get(key.Value, Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *Table) GetFloat(key Column, column string) (float64, bool, error) {
	var result float64
	err, found := table.Get(key.Value, Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *Table) GetUint(key Column, column string) (uint64, bool, error) {
	var result uint64
	err, found := table.Get(key.Value, Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *Table) GetBoolean(key Column, column string) (bool, bool, error) {
	var result bool
	err, found := table.Get(key.Value, Columns{key}, []string{column}, []interface{}{&result})
	if err != nil {
		return false, found, err
	}
	return result, found, nil
}

func (table *Table) Get(shardKey interface{}, keys Columns, columns []string, data []interface{}) (error, bool) {
	return table.Shards[table.GetShardNum(shardKey)].Get(keys, columns, data)
}
func (table *Table) Put(shardKey interface{}, values Columns) error {
	return table.Shards[table.GetShardNum(shardKey)].Put(values)
}
func (table *Table) Set(shardKey interface{}, keys Columns, values Columns) error {
	return table.Shards[table.GetShardNum(shardKey)].Set(keys, values)
}
func (table *Table) Add(shardKey interface{}, keys Columns, values Columns) error {
	return table.Shards[table.GetShardNum(shardKey)].Add(keys, values)
}
func (table *Table) Remove(shardKey interface{}, keys Columns) error {
	return table.Shards[table.GetShardNum(shardKey)].Remove(keys)
}

func (table *Table) Exec(query string, key interface{}) (sql.Result, error) {
	shardNum := table.GetShardNum(key)
	return table.Shards[shardNum].Exec(query)
}
func (table *Table) Query(query string, key interface{}) (*sql.Rows, error) {
	shardNum := table.GetShardNum(key)
	return table.Shards[shardNum].Query(query)
}

func (table *Table) ReleaseRows(rows *sql.Rows) error {
	return rows.Close()
}

func (table *Table) SingleSet(shardKey interface{}, keys Columns, column Column) error {
	return table.Set(shardKey, keys, Columns{column})
}

func (table *Table) DropUnsafe() {
	for i := 0; i < len(table.Shards); i++ {
		table.Shards[i].Drop()
	}
}

func (table *Table) GlobalExecUnsafe(query string) error {
	for i := 0; i < len(table.Shards); i++ {
		_, err := table.Shards[i].Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
