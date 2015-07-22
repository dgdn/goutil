package common

import (
	"encoding/json"
)

func LogJsonf(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
