package hparam

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type GetKey interface {
	GetKey(Key) string
}

type request struct {
	*http.Request
	err    error
	getKey func(Key) string
}

type Key string

func From(r *http.Request) *request {
	return &request{r, nil, func(key Key) string {
		return r.FormValue(key.String())
	}}
}
func FromFunc(getKey GetKey) *request {
	return &request{nil, nil, func(key Key) string {
		return getKey.GetKey(key)
	}}
}
func Optional(k string) Key {
	return Key(k + "|O")
}
func (k Key) String() string {
	if k.isOptional() {
		return strings.Split(string(k), "|")[0]
	}
	return string(k)
}
func (k Key) isOptional() bool {
	return strings.Contains(string(k), "|O")
}

func (p *request) Scan(key Key, dest interface{}) *request {
	if p.err != nil {
		return p
	}
	sval := p.getKey(key)
	if key.isOptional() == false && sval == "" {
		p.err = fmt.Errorf("missing [%v] hparam", key)
		return p
	}
	if key.isOptional() && sval == "" {
		return p
	}
	p.err = scanVal(sval, dest)
	return p
}

func (p *request) ScanM(keys []Key, dests ...interface{}) *request {

	if len(dests) == 1 {

		typ := reflect.TypeOf(dests[0]).Elem()
		if v, ok := dests[0].(ScanP); ok {
			p.err = v.ScanP(keys[0].String())
			return p
		}
		if typ.Kind() == reflect.Struct {
			return p.ScanS(keys, dests[0])
		}
	}
	if len(keys) != len(dests) {
		p.err = fmt.Errorf("keys len %v not equal to dests len %v", len(keys), len(dests))
		return p
	}
	for i := range keys {
		p.Scan(keys[i], dests[i])
	}
	return p
}

func (p *request) ScanS(keys []Key, dest interface{}) *request {
	if p.err != nil {
		return p
	}

	//for each key
	//key,dest val---> interface{}
	//scan (key interface{})
	for _, key := range keys {
		pval := getValueByKey(key, dest)
		p.Scan(key, pval)
	}
	return p

}
func getValueByKey(key Key, dest interface{}) interface{} {
	v := getFieldValueOption(key, dest)
	if v == nil {
		panic(fmt.Errorf("can not find tag[%v] in the struct type[%T]", key, dest))
	}
	return v
}

func getFieldValueOption(key Key, dest interface{}) interface{} {

	val := reflect.ValueOf(dest).Elem()
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Tag.Get("hparam") == key.String() {
			return val.Field(i).Addr().Interface()
		}
	}
	return nil
}

type ScanP interface {
	ScanP(string) error
}

func ScanToField(key Key, sval string, dest interface{}) error {
	des := getFieldValueOption(key, dest)
	if des == nil {
		return nil
	}
	return scanVal(sval, des)
}

func scanVal(sval string, dest interface{}) error {
	reflectVal := reflect.ValueOf(dest)
	kind := reflectVal.Kind()
	if kind != reflect.Ptr {
		return fmt.Errorf("dest val type is not ptr")
	}
	elem := reflectVal.Elem()
	switch elem.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		return parseInt(sval, elem)
	case reflect.Float32, reflect.Float64:
		return parseFloat(sval, elem)
	case reflect.Bool:
		return parseBool(sval, elem)
	case reflect.String:
		return parseString(sval, elem)
	// case reflect.Slice:
	// 	return parseSlice(sval, elem)
	default:
		// fmt.Println(elem.Kind())
		p, ok := reflectVal.Interface().(ScanP)
		if ok {
			return p.ScanP(sval)
		} else {
			return fmt.Errorf("not support the type %T", dest)
		}
	}
	return nil
}
func parseSlice(sval string, rval reflect.Value) error {

	switch rval.Elem().Kind() {
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		return parseIntSlice(sval, rval)
	case reflect.Float32, reflect.Float64:
		return parseFloatSlice(sval, rval)
	case reflect.String:
		return parseStringSlice(sval, rval)
	default:
		return fmt.Errorf("not support the Slice type %T", rval)
	}
}
func parseIntSlice(sval string, rval reflect.Value) error {
	for _, s := range strings.Split(sval, ",") {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		rval = reflect.Append(rval, reflect.ValueOf(v))
	}
	return nil
}
func parseStringSlice(sval string, rval reflect.Value) error {
	for _, s := range strings.Split(sval, ",") {
		rval = reflect.Append(rval, reflect.ValueOf(s))
	}
	return nil
}
func parseFloatSlice(sval string, rval reflect.Value) error {
	for _, s := range strings.Split(sval, ",") {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		rval = reflect.Append(rval, reflect.ValueOf(v))
	}
	return nil
}
func parseString(sval string, rval reflect.Value) error {
	rval.SetString(sval)
	return nil
}
func parseBool(sval string, rval reflect.Value) error {
	v, err := strconv.ParseInt(sval, 10, 64)
	if err == nil {
		if v > 0 {
			rval.SetBool(true)
			return nil
		} else {
			rval.SetBool(false)
			return nil
		}
	}
	if sval == "true" {
		rval.SetBool(true)
		return nil
	} else if sval == "false" {
		rval.SetBool(false)
		return nil
	}
	return fmt.Errorf("can not parse [%v] to bool type", sval)
}
func parseInt(sval string, rval reflect.Value) error {
	v, err := strconv.ParseInt(sval, 10, 64)
	if err != nil {
		return err
	}
	rval.SetInt(v)
	return nil
}
func parseFloat(sval string, rval reflect.Value) error {
	v, err := strconv.ParseFloat(sval, 64)
	if err != nil {
		return err
	}
	rval.SetFloat(v)
	return nil
}
func (p *request) Error() error {
	return p.err
}
