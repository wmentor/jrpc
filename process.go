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

func Process(in io.Reader, out io.Writer) error {
	return global.Process(in, out)
}

func (c JRPC) Process(in io.Reader, out io.Writer) error {

	input, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	input = bytes.TrimSpace(input)

	if len(input) > 0 {
		if input[0] == '{' {

			var rec Request

			if err := json.Unmarshal(input, &rec); err != nil {
				return ErrInvalidRequest
			}

			if resp := c.exec(&rec); resp != nil {
				return json.NewEncoder(out).Encode(resp)
			}

			return nil

		} else if input[0] == '[' {

			recs := make([]Request, 0)
			resp := make([]interface{}, 0, 10)

			if err := json.Unmarshal(input, &recs); err != nil {
				return ErrInvalidRequest
			}

			for _, rec := range recs {
				if r := c.exec(&rec); r != nil {
					resp = append(resp, r)
				}
			}

			return json.NewEncoder(out).Encode(resp)
		}
	}

	return ErrInvalidRequest
}
