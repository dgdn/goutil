package common

import (
	"com.dy.rcp/dbMgr"
	_ "com.dy.rcp/testconf"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	db := dbMgr.DbConn()

	d1 := StringSlice{"2", "33"}
	var d2 StringSlice

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

	if d2[0] != "2" {
		t.Errorf("d[0] not equal to %v", 2)
	}
}

func TestJsonMarsh(t *testing.T) {
	d1 := StringSlice{"2", "33"}
	d, err := json.Marshal(d1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(d))
	if string(d) != `["2","33"]` {
		t.Error("encoded string not equal")
	}

	var d2 StringSlice
	err = json.Unmarshal(d, &d2)
	if err != nil {
		t.Error(err)
	}
	if d2[0] != "2" {
		t.Errorf("d2[0] not equal to %v", 2)
	}
}
