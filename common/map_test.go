package common

import (
	"com.dy.rcp/dbMgr"
	_ "com.dy.rcp/testconf"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestMapScan(t *testing.T) {
	db := dbMgr.DbConn()

	type testStruct struct {
		Test string
	}

	d1 := Map{
		"test": "hehe",
		"val":  testStruct{"kao"},
	}
	var d2 Map

	db.Exec(`delete from ucs_usr where tid=?`, 99999)

	_, err := db.Exec(`insert into ucs_usr(tid,usr,pwd,time,status,add1)values(?,?,?,?,?,?)`,
		99999, "hehe", "er", time.Now(), "hehe", d1)
	if err != nil {
		t.Error(err)
	}

	err = db.QueryRow(`select add1 from ucs_usr where tid=?`, 99999).Scan(&d2)
	if err != nil {
		t.Error(err)
	}

	fmt.Print(d2)

	if d2["test"] != "hehe" {
		t.Errorf(`d["test"] not equal to %v`, "hehe")
	}
}

func TestMapJsonMarsh(t *testing.T) {
	d1 := Map{"test": "hehe"}
	d, err := json.Marshal(d1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(d))
	if string(d) != `{"test":"hehe"}` {
		t.Error("encoded string not equal")
	}

	var d2 Map
	err = json.Unmarshal(d, &d2)
	if err != nil {
		t.Error(err)
	}
	if d2["test"] != "hehe" {
		t.Errorf(`d2["test"] not equal to %v`, "hehe")
	}
}

func TestValidator(t *testing.T) {
	err := ValidateParam(
		MustBePositive(-1, "orgId"),
		MustBePositive(-1, "auditor"),
	)
	fmt.Println(err)
}
