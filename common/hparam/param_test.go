package hparam

import (
	"net/http"
	"strings"
	"testing"
)

func TestParsehparamOptional(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	var a int64
	if err := From(r).Scan(Optional("a"), &a).Error(); err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("expect a equal %v got %v", 1, a)
	}

	var c int64
	if err := From(r).Scan(Optional("c"), &c).Error(); err != nil {
		t.Error(err)
	}
	if c != 0 {
		t.Errorf("expect c equal %v got %v", 0, c)
	}

}

func TestParsehparamInt(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1&b=ac", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	var a int64
	if err := From(r).Scan("a", &a).Error(); err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("expect a equal %v got %v", 1, a)
	}

	if err := From(r).Scan("b", &a).Error(); err == nil {
		t.Error(err)
	}
}

func TestParsehparamBool(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1&b=0&c=true&d=false&e=k", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	var a bool
	if err := From(r).Scan("a", &a).Error(); err != nil {
		t.Error(err)
	}
	if a != true {
		t.Errorf("expect a equal %v got %v", true, a)
	}

	if err := From(r).Scan("b", &a).Error(); err != nil {
		t.Error(err)
	}
	if a != false {
		t.Errorf("expect b equal %v got %v", false, a)
	}

	if err := From(r).Scan("c", &a).Error(); err != nil {
		t.Error(err)
	}
	if a != true {
		t.Errorf("expect c equal %v got %v", true, a)
	}

	if err := From(r).Scan("d", &a).Error(); err != nil {
		t.Error(err)
	}
	if a != false {
		t.Errorf("expect d equal %v got %v", false, a)
	}

	if err := From(r).Scan("e", &a).Error(); err == nil {
		t.Error("expect err")
	}

}

func TestParseString(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1&b=ac", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	var a string
	if err := From(r).Scan("a", &a).Error(); err != nil {
		t.Error(err)
	}
	if a != "1" {
		t.Errorf("expect a equal %v got %v", "1", a)
	}
}

type ssType []string

func (ss *ssType) ScanP(in string) error {
	*ss = strings.Split(in, ",")
	return nil
}
func TestParseCustomType(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1,2&b=ac", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	var ss ssType
	if err := From(r).Scan("a", &ss).Error(); err != nil {
		t.Error(err)
	}
	if len(ss) != 2 {
		t.Errorf("epect len ss %v got %v", 2, len(ss))
	}
	if ss[0] != "1" || ss[1] != "2" {
		t.Errorf("expect ss[0] ss[1] %v %v got %v %v", "1", "2", ss[0], ss[1])
	}

}

type aaType struct {
	A string
}

func (aa *aaType) ScanP(in string) error {
	aa.A = in
	return nil
}
func TestParseArray(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1,2,3&b=ac", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	var is []int
	if err := From(r).Scan("a", &is).Error(); err != nil {
		t.Error(err)
	}
	if len(is) != 3 {
		t.Errorf("expect len [is] %v but %v", 3, len(is))
	}
	if is[0] != 1 && is[1] != 2 {
		t.Errorf("expect [is] %v but %v", []int{1, 2}, is)
	}
}

func TestParseCustomTypeStruct(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1&b=ac", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	var ss aaType
	if err := From(r).Scan("a", &ss).Error(); err != nil {
		t.Error(err)
	}
	if ss.A != "1" {
		t.Errorf("expect A %v but got %v", "1", ss.A)
	}
}

func TestParseNotSuppotType(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1,2&b=ac", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	var ss []string
	if err := From(r).Scan("a", &ss).Error(); err == nil {
		t.Error(err)
	}
}

func TestScanStruct(t *testing.T) {
	r, err := http.NewRequest("GET", "http://hello.com?a=1&b=ac&c=true", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	type T struct {
		A int64  `hparam:"a"`
		B string `hparam:"b"`
		C bool   `hparam:"c"`
	}
	var s T
	if err := From(r).ScanS([]Key{"a|O", "b", "c"}, &s).Error(); err != nil {
		t.Error(err)
	}
	if s.A != 1 {
		t.Errorf("expect s.A to be %v got %v", 1, s.A)
	}
	if s.B != "ac" {
		t.Errorf("expect s.B to be %v got %v", "ac", s.B)
	}
	if s.C != true {
		t.Errorf("expect s.C to be %v got %v", true, s.C)
	}

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("should panic")
			}
		}()
		if err := From(r).ScanS([]Key{"no"}, &s).Error(); err != nil {
			t.Error(err)
		}
	}()

}
