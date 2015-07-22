package common

import (
	"com.dy.rcp/common/hparam"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	_common "org.cny.uas/common"
	"reflect"
)

type getkey struct {
	hs *routing.HTTPSession
}

func newGetKey(hs *routing.HTTPSession) getkey {
	return getkey{hs}
}

func (g getkey) GetKey(key hparam.Key) string {
	if key.String() == "{uid}" {
		return fmt.Sprint(g.hs.IntVal("uid"))
	}
	return g.hs.R.FormValue(key.String())
}

func MakeHttpApi(fn interface{}, params ...hparam.Key) routing.HandleFunc {

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic("first param must be func")
	}

	return wrapHandle(func(hs *routing.HTTPSession) []reflect.Value {
		if hs.R.Method == "POST" {
			return post(fn, hs)
		}
		return get(fn, hs, params)
	})

}

func MakeHttpRedirectApi(fn interface{}, params ...hparam.Key) routing.HandleFunc {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic("first param must be func")
	}
	return wrapRedirectHandle(func(hs *routing.HTTPSession) []reflect.Value {
		if hs.R.Method == "POST" {
			return post(fn, hs)
		}
		return get(fn, hs, params)
	})

}

func get(fn interface{}, hs *routing.HTTPSession, scanKeys []hparam.Key) []reflect.Value {
	typFn := reflect.TypeOf(fn)
	var callParams []reflect.Value
	var scanValues []interface{}

	//none-ptr
	for i := 0; i < typFn.NumIn(); i++ {
		if typFn.In(i).Kind() == reflect.Ptr {
			var val = reflect.New(typFn.In(i).Elem())
			scanValues = append(scanValues, val.Interface())
			callParams = append(callParams, val)
		} else {
			var val = reflect.New(typFn.In(i))
			scanValues = append(scanValues, val.Interface())
			callParams = append(callParams, val.Elem())
		}

	}

	fmt.Printf("%v %+v \n", scanKeys, scanValues)
	if err := hparam.FromFunc(newGetKey(hs)).ScanM(scanKeys, scanValues...).Error(); err != nil {
		return []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(err)}
	}
	valFn := reflect.ValueOf(fn)
	return valFn.Call(callParams)
}

func post(fn interface{}, hs *routing.HTTPSession) []reflect.Value {
	typFn := reflect.TypeOf(fn)
	if typFn.NumIn() != 1 {
		panic("fn handle param num must be 1")
	}
	var callParam reflect.Value
	var destParam reflect.Value
	if typFn.In(0).Kind() == reflect.Ptr {
		destParam = reflect.New(typFn.In(0).Elem())
		callParam = destParam
	} else {
		destParam = reflect.New(typFn.In(0))
		callParam = destParam.Elem()
	}

	if err := JsonUnmarshalBody(hs, destParam.Interface()); err != nil {
		return []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(err)}
	}

	if err := hparam.ScanToField("{uid}", newGetKey(hs).GetKey("{uid}"), destParam.Interface()); err != nil {
		return []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(err)}
	}

	valFn := reflect.ValueOf(fn)
	return valFn.Call([]reflect.Value{callParam})
}

func wrapRedirectHandle(fn func(*routing.HTTPSession) []reflect.Value) routing.HandleFunc {
	return func(hs *routing.HTTPSession) routing.HResult {
		outs := fn(hs)
		if len(outs) != 2 {
			panic("redirect handle api return count should be 2")
		}
		if err, ok := outs[1].Interface().(error); ok && err != nil {
			log.E("err:%v", err.Error())
			return hs.MsgResE(1, err.Error())
		}
		hs.Redirect(outs[0].Interface().(string))
		return _common.MsgRes(hs, outs[0].Interface().(string))
	}
}

func wrapHandle(fn func(*routing.HTTPSession) []reflect.Value) routing.HandleFunc {
	return func(hs *routing.HTTPSession) routing.HResult {
		outs := fn(hs)
		if len(outs) == 2 {
			return SimpleHandle(func(hs *routing.HTTPSession) (interface{}, error) {
				var err error
				if outs[1].Interface() != nil {
					err = outs[1].Interface().(error)
				}
				return outs[0].Interface(), err
			})(hs)
		} else if len(outs) == 3 {
			return SimpleListHandle(func(hs *routing.HTTPSession) (interface{}, int64, error) {
				var err error
				if outs[2].Interface() != nil {
					err = outs[2].Interface().(error)
				}
				return outs[0].Interface(), outs[1].Int(), err
			})(hs)
		} else {
			panic("return count should be 2 or 3")
		}
		return routing.HRES_RETURN
	}
}
