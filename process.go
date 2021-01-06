package jrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

var (
	ErrInvalidRequest error = errors.New("invalid request")
)

func Process(in io.Reader) ([]byte, error) {
	return global.Process(in)
}

func (c JRPC) Process(in io.Reader) ([]byte, error) {

	input, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	input = bytes.TrimSpace(input)

	if len(input) > 0 {
		if input[0] == '{' {

			var rec Request

			if err := json.Unmarshal(input, &rec); err != nil {
				return nil, ErrInvalidRequest
			}

			if resp := c.exec(&rec); resp != nil {
				return json.Marshal(resp)
			}

			return nil, nil

		} else if input[0] == '[' {

			recs := make([]Request, 0)
			resp := make([]interface{}, 0, 10)

			if err := json.Unmarshal(input, &recs); err != nil {
				return nil, ErrInvalidRequest
			}

			for _, rec := range recs {
				if r := c.exec(&rec); r != nil {
					resp = append(resp, r)
				}
			}

			return json.Marshal(resp)
		}
	}

	return nil, ErrInvalidRequest
}
