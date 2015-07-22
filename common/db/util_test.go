package dbMgr

import (
	"fmt"
	"testing"
)

type T struct {
	A string `m2s:"a,autoinc"`
	B string `m2s:"b"`
}

func (v *T) GetTableName() string {
	return "rcp_test_table"
}
func TestParseField(t *testing.T) {

	tv := T{"a", "b"}
	fields := parseToDbfield(&tv)
	for _, f := range fields {
		fmt.Printf("%+v", f)
	}
	Insert(&tv)
}
