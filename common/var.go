package common

import (
	"errors"
	"fmt"
	"github.com/Centny/gwf/log"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	DB_ERROR       = "系统数据异常"
	ARG_ERROR      = "输入参数错误"
	DATA_NOT_FOUND = "请求的数据未找到"
)

type ObjectId struct {
	bson.ObjectId
}

func (o *ObjectId) ScanP(s string) error {
	if !bson.IsObjectIdHex(s) {
		return fmt.Errorf("id format err")
	}
	o.ObjectId = bson.ObjectIdHex(s)
	return nil
}

func LogErr(err error, reason string) error {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		log.E("%v:%v %v:%v", file, line, reason, err.Error())

	} else {
		log.E(reason)
	}
	return errors.New(reason)
}

func LogNErr(err error) error {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		log.E("%v:%v %v", file, line, err.Error())

	}
	return err
}

func GetMapInt(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		return v.(int64)
	}
	return -1
}

func GetMapFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		return v.(float64)
	}
	return -1
}

func GetMapVal(m map[string]interface{}, key string) interface{} {
	if v, ok := m[key]; ok {
		return v
	}
	return nil
}

var NowTimeString = func() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type RangeInt struct {
	From  int64
	To    int64
	Valid bool
}
type RangeFloat struct {
	From  float64
	To    float64
	Valid bool
}

func MakeRangeFloat(s string) (RangeFloat, error) {
	if s == "" {
		return RangeFloat{}, nil
	}

	rf := RangeFloat{}
	var err error = fmt.Errorf("range float format wrong")
	sa := strings.Split(s, ",")
	if len(sa) != 2 {
		return rf, err
	}
	rf.From, err = strconv.ParseFloat(sa[0], 64)
	if err != nil {
		return rf, err
	}
	rf.To, err = strconv.ParseFloat(sa[1], 64)
	if err != nil {
		return rf, err
	}
	rf.Valid = true
	return rf, nil
}

func MakeRangeInt(s string) (RangeInt, error) {
	if s == "" {
		return RangeInt{}, nil
	}

	ri := RangeInt{}
	var err error = fmt.Errorf("range int format wrong")
	sa := strings.Split(s, ",")
	if len(sa) != 2 {
		return ri, err
	}
	ri.From, err = strconv.ParseInt(sa[0], 10, 64)
	if err != nil {
		return ri, err
	}
	ri.To, err = strconv.ParseInt(sa[1], 10, 64)
	if err != nil {
		return ri, err
	}
	ri.Valid = true
	return ri, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type NullInt64 struct {
	Int64 int64
	Valid bool
}

func (m *NullInt64) UnmarshalJSON(data []byte) (err error) {
	if m == nil {
		return errors.New("NullInt64: UnmarshalJSON on nil pointer")
	}
	str := string(data)
	m.Int64, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		m.Valid = false
		return
	}
	m.Valid = true
	return
}

type NullString struct {
	String string
	Valid  bool
}
type NullBool struct {
	Bool  bool
	Valid bool
}

func NewNullString(s string) NullString {
	if s == "" {
		return NullString{s, false}
	}
	return NullString{s, true}
}

func NullBoolFromString(s string) NullBool {
	if s == "true" {
		return NullBool{true, true}
	} else if s == "false" {
		return NullBool{false, true}
	}
	return NullBool{false, false}
}

// type JsonTime time.Time

// // type JsonDuration time.Duration

// func (t JsonTime) MarshalJSON() ([]byte, error) {
// 	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
// 	return []byte(stamp), nil
// }

func CheckRun(t ...interface{}) {
	err, found := t[len(t)-1].(error)
	if found {
		panic(err)
	}
}
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
