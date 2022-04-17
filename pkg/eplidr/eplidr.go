package eplidr

import (
	"fmt"
	"strconv"
)

type Column struct {
	Key   string
	Value interface{}
}

func (column *Column) GetStringValue() string {
	return value(column.Value)
}

func value(i interface{}) string {
	switch v := i.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case []interface{}: // Serialize s
		return fmt.Sprintf("'%v'", v[0])
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'E', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type Columns []Column

func KeysToQuery(keys Columns) string {
	if len(keys) == 0 {
		return ""
	}
	query := "WHERE"
	for i := 0; i < len(keys); i++ {
		if i == len(keys)-1 {
			query += fmt.Sprintf("`%s`=%s", keys[i].Key, keys[i].GetStringValue())
		} else {
			query += fmt.Sprintf("`%s`=%s,", keys[i].Key, keys[i].GetStringValue())
		}
	}
	return query
}

func ColumnNamesToQuery(names ...string) string {
	result := ""
	for i := 0; i < len(names); i++ {
		result += " `" + names[i] + "`"
	}
	return result
}

func PlainToColumns(keys []string, values []interface{}) Columns {
	columns := make(Columns, len(keys))
	for i := 0; i < len(keys); i++ {
		columns[i] = Column{
			Key:   keys[i],
			Value: values[i],
		}
	}
	return columns
}

var shiftTableIterator int

func init() {
	shiftTableIterator = 1
}

func StandardGetShardFunc(key interface{}) uint {
	return fnv32(fmt.Sprintf("%v", key))
}
func fnv32(key string) uint {
	hash := uint(2166136261)
	const prime32 = uint(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint(key[i])
	}
	return hash
}
