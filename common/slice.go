package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type StringSlice []string

func (n *StringSlice) Scan(value interface{}) error {

	switch value.(type) {
	case string:
		return json.Unmarshal([]byte(value.(string)), n)
	case nil:
		n = &StringSlice{}
		return nil
	case []byte:
		return json.Unmarshal(value.([]byte), n)
	}
	return fmt.Errorf("no support the value type :%T", value)

}

// Value implements the driver Valuer interface.
func (n StringSlice) Value() (driver.Value, error) {
	return json.Marshal(n)
}

type Int64Slice []int64

func (n *Int64Slice) Scan(value interface{}) error {

	switch value.(type) {
	case string:
		return json.Unmarshal([]byte(value.(string)), n)
	case nil:
		n = &Int64Slice{}
		return nil
	case []byte:
		return json.Unmarshal(value.([]byte), n)
	}
	return fmt.Errorf("no support the value type :%T", value)

}

func I64In(ins []int64, e int64) bool {
	for _, v := range ins {
		if e == v {
			return true
		}
	}
	return false
}

func (n Int64Slice) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func Int64SliceToJsonStr(slice []int64) string {
	a, _ := json.Marshal(slice)
	return string(a)
}

func Int64SliceToCommaStr(slice []int64) string {
	strs := []string{}
	for _, i := range slice {
		strs = append(strs, strconv.FormatInt(i, 10))
	}
	return strings.Join(strs, ",")
}

func CommaStrToInt64Slice(a string) ([]int64, error) {
	is := []int64{}
	strs := strings.Split(a, ",")
	for _, s := range strs {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		is = append(is, i)
	}
	return is, nil

}
