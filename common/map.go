package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Map map[string]interface{}

func (n *Map) Scan(value interface{}) error {

	switch value.(type) {
	case string:
		return json.Unmarshal([]byte(value.(string)), n)
	case nil:
		n = &Map{}
		return nil
	case []byte:
		return json.Unmarshal(value.([]byte), n)
	}
	return fmt.Errorf("no support the value type :%T", value)

}

// Value implements the driver Valuer interface.
func (n Map) Value() (driver.Value, error) {
	return json.Marshal(n)
}
