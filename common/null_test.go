package common

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshNullInt64(t *testing.T) {
	type T struct {
		V NullInt64
	}
	var ni T
	err := json.Unmarshal([]byte(`{}`), &ni)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal([]byte(`{"V":"111"}`), &ni)
	if err == nil {
		t.Error(err)
	}

	err = json.Unmarshal([]byte(`{"V":111}`), &ni)
	if err != nil {
		t.Error(err)
	}
	if ni.V.Int64 != 111 {
		t.Error("expect %v got %v", 111, ni.V.Int64)
	}
	fmt.Printf("%+v", ni)
}
