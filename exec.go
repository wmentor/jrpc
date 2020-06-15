package jrpc

import (
	"encoding/json"
	"reflect"
)

var callback map[string]interface{} = map[string]interface{}{}

func RegMethod(method string, fn interface{}) {
	callback[method] = fn
}

func MakeError(rec *Request, code int64, message string) *Response {
	return &Response{Id: rec.Id, JsonRPC: "2.0", Error: &Error{Code: code, Message: message}}
}

func Exec(rec *Request) (resp *Response) {

	defer func() {

		if r := recover(); r != nil {
			resp = MakeError(rec, -32600, "Invalid Request")
		}

	}()

	if rec.JsonRPC != "2.0" {
		if rec.Id != nil {
			return MakeError(rec, -32600, "Invalid Request")
		}
		return nil
	}

	fn, has_fn := callback[rec.Method]
	if !has_fn {
		if rec.Id != nil {
			return MakeError(rec, -32601, "Method not found")
		}
		return nil
	}

	if reflect.TypeOf(fn).NumIn() != 1 {
		if rec.Id != nil {
			return MakeError(rec, -32603, "Internal error")
		}
		return nil
	}

	arg_type := reflect.TypeOf(fn).In(0)
	arg := reflect.New(arg_type)
	arg_i := arg.Interface()

	data, _ := rec.Params.MarshalJSON()

	if err := json.Unmarshal(data, &arg_i); err != nil {
		if rec.Id != nil {
			return MakeError(rec, -32602, "Invalid params")
		}
		return nil
	}

	rets := reflect.ValueOf(fn).Call([]reflect.Value{arg.Elem()})

	if rec.Id == nil {
		return nil
	}

	if !rets[1].IsNil() {
		return &Response{Id: rec.Id, JsonRPC: rec.JsonRPC, Error: rets[1].Interface()}
	}

	return &Response{Id: rec.Id, JsonRPC: rec.JsonRPC, Result: rets[0].Interface()}
}
