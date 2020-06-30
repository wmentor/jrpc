package jrpc

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

type Request struct {
	Method  string           `json:"method"`
	JsonRPC string           `json:"jsonrpc"`
	Params  *json.RawMessage `json:"params"`
	Id      *json.RawMessage `json:"id"`
}

type Response struct {
	Id      *json.RawMessage `json:"id"`
	JsonRPC string           `json:"jsonrpc"`
	Result  interface{}      `json:"result"`
}

type ErrResponse struct {
	Id      *json.RawMessage `json:"id"`
	JsonRPC string           `json:"jsonrpc"`
	Error   interface{}      `json:"error"`
}
