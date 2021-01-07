package jrpc

import (
	"bytes"
	"strings"
	"testing"
)

func TestJRPC(t *testing.T) {

	RegMethod("user.hello", func(name string) (string, *Error) {
		return "Hello, " + name, nil
	})

	RegMethod("echo", func(data map[string]interface{}) (interface{}, *Error) {
		return data, nil
	})

	tP := func(in string, out string) {
		buf := bytes.NewBuffer(nil)
		if err := Process(strings.NewReader(in), buf); err != nil {
			t.Fatal(err)
		}
		if buf.String() != out {
			t.Fatalf("failed: %s ret: %s wait: %s", in, buf.String(), out)
		}
	}

	tP(`{"jsonrpc":"2.0","id":1,"method":"user.hello","params":"Mike"}`, `{"id":1,"jsonrpc":"2.0","result":"Hello, Mike"}
`)

	tP(`{"jsonrpc":"2.0","method":"user.hello","params":"Mike"}`, ``)

	tP(`{"jsonrpc":"2.0","id":1,"method":"user.add","params":"Mike"}`, `{"id":1,"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"}}
`)

	tP(`{"jsonrpc":"2.0","id":"2134-42341","method":"echo","params":{"name":"Mike"}}`, `{"id":"2134-42341","jsonrpc":"2.0","result":{"name":"Mike"}}
`)
}
