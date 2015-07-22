package common

import (
	"fmt"
	"testing"
)

func TestSqlBuilder(t *testing.T) {
	sb := NewSqlBuilder().Select("cid,name,id").
		From("rcp_section s left join rcp_course c on c.tid=s.tid").
		Where("and cid=?", 1).
		Where("and uid=?", 1).
		Where("and status='N'", nil).
		Limit(1, 10).
		OrderBy("tid asc")
	fmt.Println(sb.SelectSql())
	fmt.Println(sb.CountSql())
	fmt.Println(sb.Args)
	expected := `select cid,name,id from rcp_section s left join rcp_course c on c.tid=s.tid where 1=1 and cid=? and uid=? and status='N' order by tid asc limit 0,10`
	if sb.SelectSql() != expected {
		t.Errorf("expected %v \ngot %v", expected, sb.SelectSql())
	}
	expected2 := `select count(*) from rcp_section s left join rcp_course c on c.tid=s.tid where 1=1 and cid=? and uid=? and status='N'`
	if sb.CountSql() != expected2 {
		t.Errorf("expected %v \n got %v", expected2, sb.CountSql())
	}
	if sb.Args[0] != 1 || sb.Args[1] != 1 {
		t.Errorf("args fail %v", sb.Args)
	}
}
