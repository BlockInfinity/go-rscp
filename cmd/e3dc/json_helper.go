package main

import (
	"bytes"
	"encoding/json"

	"github.com/BlockInfinity/go-rscp/rscp"
)

func isJSONEmpty(d []byte) bool {
	x := bytes.TrimLeft(d, " \t\r\n")
	return len(x) == 0
}

func isJSONArray(d []byte) bool {
	x := bytes.TrimLeft(d, " \t\r\n")
	return len(x) > 0 && x[0] == '['
}

func isJSONString(d []byte) bool {
	x := bytes.TrimLeft(d, " \t\r\n")
	return len(x) > 0 && x[0] == '"'
}

func isJSONNumber(d []byte) bool {
	x := bytes.TrimLeft(d, " \t\r\n")
	return len(x) > 0 && x[0] >= 48 && x[0] <= 57
}

func isJSONDataType(d []byte) (bool, *rscp.DataType) {
	if isJSONString(d) {
		s := new(rscp.DataType)
		if err := json.Unmarshal(d, s); err == nil && s.IsADataType() {
			return true, s
		}
	}
	return false, nil
}
