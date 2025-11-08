package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func (h Headers) Get(key string) (string, bool) {
	val, ok := h[strings.ToLower(key)]
	if !ok {
		return "", ok
	}
	return val, ok
}

func (h Headers) Set(key string, value string) {
	h[strings.ToLower(string(key))] = value
}

func (h Headers) Delete(key string) {
	delete(h, key)
}

var seperator = []byte("\r\n")

var ErrWrongFieldFormat = errors.New("wrong headers format provided")
var ErrWrongFieldKeyFormat = errors.New("wrong headers key format provided")
var ErrUnsupportedCharacter = errors.New("unsupported character in field name")
var ErrWrongFieldValueFormat = errors.New("wrong headers value format provided")

func (h Headers) validateFieldName(fieldName []byte) bool {
	for _, ch := range fieldName {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '!' || ch == '#' || ch == '$' || ch == '&' || ch == '*' || ch == '-' || ch == '%' || ch == '+' || ch == '.' || ch == '^' || ch == '_' || ch == '`' || ch == '|' || ch == '~' || ch == '\'' {
			continue
		} else {
			return false
		}
	}
	return true
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	n := 0
	done := false
	var err error = nil
	for {
		idx := bytes.Index(data, seperator)
		if idx == 0 {
			n += len(seperator)
			done = true
			break
		}
		if idx == -1 {
			break
		}
		n += idx + len(seperator)
		fieldLine := data[:idx]
		data = data[idx+len(seperator):]
		fields := bytes.SplitN(fieldLine, []byte(":"), 2)
		if len(fields) != 2 {
			n = 0
			err = ErrWrongFieldFormat
			break
		}
		fieldKey := fields[0]
		fieldValue := fields[1]
		fieldValue = bytes.Trim(fieldValue, " ")
		if bytes.ContainsAny(fieldKey, " ") {
			n = 0
			err = ErrWrongFieldFormat
			break
		}
		if !h.validateFieldName(fieldKey) {
			n = 0
			err = ErrWrongFieldValueFormat
			break
		}
		if val, ok := h[strings.ToLower(string(fieldKey))]; ok {
			h[strings.ToLower(string(fieldKey))] = val + fmt.Sprintf(",%s", fieldValue)
		} else {
			h[strings.ToLower(string(fieldKey))] = string(fieldValue)
		}
	}
	return n, done, err
}

func NewHeaders() Headers {
	return make(Headers, 0)
}
