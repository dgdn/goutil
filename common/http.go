package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"io/ioutil"
	_common "org.cny.uas/common"
)

type myReader struct {
	*bytes.Buffer
}

// So that it implements the io.ReadCloser interface
func (m myReader) Close() error { return nil }

func JsonUnmarshalBody(hs *routing.HTTPSession, v interface{}) error {
	buf, _ := ioutil.ReadAll(hs.R.Body)
	rdr1 := myReader{bytes.NewBuffer(buf)}
	rdr2 := myReader{bytes.NewBuffer(buf)}

	decoder := json.NewDecoder(rdr1)
	err := decoder.Decode(v)
	if err != nil {
		log.E("parse request body fail:%v", err.Error())
		return fmt.Errorf("请求参数格式错误:%v", err)
	}
	hs.R.Body = rdr2
	return nil
}

func SimpleHandle(fn func(*routing.HTTPSession) (interface{}, error)) routing.HandleFunc {
	return func(hs *routing.HTTPSession) routing.HResult {
		rst, err := fn(hs)
		if err != nil {
			log.E("err:%v", err.Error())
			return hs.MsgResE(1, err.Error())
		}
		return _common.MsgRes(hs, rst)
	}
}

func SimpleListHandle(fn func(*routing.HTTPSession) (interface{}, int64, error)) routing.HandleFunc {
	return func(hs *routing.HTTPSession) routing.HResult {
		list, total, err := fn(hs)
		if err != nil {
			log.E("err:%v", err.Error())
			return hs.MsgResE(1, err.Error())
		}
		return _common.MsgRes(hs, map[string]interface{}{
			"total": total,
			"list":  list,
		})
	}
}
func MakeDataAccessFilter(f func(*routing.HTTPSession) (bool, error)) routing.HandleFunc {
	return func(hs *routing.HTTPSession) routing.HResult {
		pass, err := f(hs)
		if err != nil {
			log.E("检验当前用户数据权限失败-%v", err.Error())
			return _common.MsgResE(hs, 1, "检验当前用户数据权限失败")
		}
		if !pass {
			return _common.MsgResE(hs, 1, "当前用户没权限访问该数据")
		}
		return routing.HRES_CONTINUE
	}
}
