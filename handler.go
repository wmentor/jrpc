package jrpc

import (
	"bytes"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidRequest error = errors.New("invalid request")
)

func Handle(in []byte) ([]byte, error) {

	in = bytes.TrimSpace(in)

	if len(in) > 0 {
		if in[0] == '{' {

			var rec Request

			if err := json.Unmarshal(in, &rec); err != nil {
				return nil, ErrInvalidRequest
			}

			if resp := Exec(&rec); resp != nil {
				out, _ := json.Marshal(resp)
				return out, nil
			}

			return nil, nil

		} else if in[0] == '[' {

			recs := make([]Request, 0)
			resp := make([]*Response, 0, 10)

			if err := json.Unmarshal(in, &recs); err != nil {
				return nil, ErrInvalidRequest
			}

			for _, rec := range recs {

				if r := Exec(&rec); r != nil {
					resp = append(resp, r)
				}
			}

			out, _ := json.Marshal(resp)
			return out, nil
		}
	}

	return nil, ErrInvalidRequest
}
