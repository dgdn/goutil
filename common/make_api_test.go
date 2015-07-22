package common

import (
	"com.dy.rcp/testconf"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go/parser"
	"go/token"
	"log"
	"os"
	"testing"
)

type ssType struct {
	A int64  `hparam:"a"`
	B string `hparam:"b"`
	C int64  `hparam:"{uid}"`
}

func TestMakeApi(t *testing.T) {
	Convey("test make api", t, func() {
		Convey("case 1", func() {
			hs, rs := testconf.HsBuilder("GET", "http://apple.com?id=1", 100, "hello", "")
			MakeHttpApi(func(id int64) (int64, error) {
				return id, nil
			}, "id")(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 0)
			So(parse["data"], ShouldEqual, 1)
		})
		Convey("case 2", func() {
			hs, rs := testconf.HsBuilder("GET", "http://apple.com?id=a1", 100, "hello", "")
			MakeHttpApi(func(id int64) (int64, error) {
				return id, nil
			}, "id")(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 1)
		})
		Convey("case 3", func() {
			hs, rs := testconf.HsBuilder("GET", "http://apple.com?id=", 100, "hello", "")
			MakeHttpApi(func(id int64) (int64, error) {
				return id, nil
			}, "id|O")(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 0)
			So(parse["data"], ShouldEqual, 0)
		})
		Convey("case 3.1", func() {
			hs, rs := testconf.HsBuilder("GET", "http://apple.com?id=", 100, "hello", "")
			MakeHttpApi(func(id int64) (int64, error) {
				return id, fmt.Errorf("bs err")
			}, "id|O")(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 1)
			So(parse["msg"], ShouldEqual, "bs err")
		})
		Convey("case 4", func() {
			hs, rs := testconf.HsBuilder("GET", "http://apple.com?id=1", 100, "hello", "")
			MakeHttpApi(func(id, uid int64) (int64, error) {
				return id + uid, nil
			}, "id", "{uid}")(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 0)
			So(parse["data"], ShouldEqual, 101)
		})
		Convey("case 5", func() {
			hs, rs := testconf.HsBuilder("GET", "http://apple.com?a=100&b=1", 100, "hello", "")
			MakeHttpApi(func(ss ssType) (string, error) {
				fmt.Println(ss)
				return fmt.Sprintf("%v%v%v", ss.A, ss.B, ss.C), nil
			}, "a", "b", "{uid}")(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 0)
			So(parse["data"], ShouldEqual, "1001100")
		})
		Convey("case 5.1", func() {
			hs, rs := testconf.HsBuilder("GET", "http://apple.com?a=100&b=1", 100, "hello", "")
			MakeHttpApi(func(ss *ssType) (string, error) {
				fmt.Println(ss)
				return fmt.Sprintf("%v%v%v", ss.A, ss.B, ss.C), nil
			}, "a", "b", "{uid}")(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 0)
			So(parse["data"], ShouldEqual, "1001100")
		})
		Convey("case 6", func() {
			hs, rs := testconf.HsBuilder("POST", "http://apple.com", 100, "hello", `{"a":1,"b":"1"}`)
			MakeHttpApi(func(ss ssType) (string, error) {
				fmt.Printf("ss %+v", ss)
				return fmt.Sprintf("%v%v%v", ss.A, ss.B, ss.C), nil
			})(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 0)
			So(parse["data"], ShouldEqual, "11100")
		})
		Convey("case 7", func() {
			hs, rs := testconf.HsBuilder("POST", "http://apple.com", 100, "hello", `{"a":1,"b":"1"}`)
			MakeHttpApi(func(ss *ssType) (string, error) {
				fmt.Printf("ss %+v", ss)
				return fmt.Sprintf("%v%v%v", ss.A, ss.B, ss.C), nil
			})(hs)
			var parse = make(map[string]interface{})
			json.Unmarshal(*rs, &parse)
			fmt.Println(parse)
			So(parse["code"], ShouldEqual, 0)
			So(parse["data"], ShouldEqual, "11100")
		})
	})
}

//@api1
//@doc1
func T1() {

}

//@api2
//@doc2
func T2() {

}

func TestParseComment(t *testing.T) {
	fileSet := token.NewFileSet()
	log.Println(os.Getenv("PWD"))
	fileTree, err := parser.ParseFile(fileSet, "./make_api_test.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Can not parse general API information: %v\n", err)
	}
	if fileTree.Comments != nil {
		for _, comment := range fileTree.Comments {
			log.Println("find comments %v", comment.Text())
		}
	} else {
		log.Fatalf("comments is nil")
	}

}
