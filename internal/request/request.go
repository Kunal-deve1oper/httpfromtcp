package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/Kunal-deve1oper/httpfromtcp/internal/headers"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       parserState
	Headers     headers.Headers
	Body        []byte
}

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
)

func (r *RequestLine) ValidateMethod(method string) bool {
	flag, _ := regexp.Match("^[A-Z]+$", []byte(method))
	return flag
}

func (r *RequestLine) ValidateHTTP() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func (r *Request) Done() bool {
	return r.State == StateDone
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		currentData := data[read:]
		switch r.State {
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.State = StateHeaders
		case StateHeaders:
			n, status, err := r.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			if status {
				r.State = StateBody
			}
		case StateBody:
			contentLen, ok := r.Headers.Get("Content-Length")
			if !ok {
				r.State = StateDone
				break outer
			}
			r.Body = append(r.Body, currentData...)
			read += len(currentData)
			bodyLen, err := strconv.Atoi(contentLen)
			if err != nil {
				return 0, err
			}
			if len(r.Body) == bodyLen {
				r.State = StateDone
				break outer
			}
			if len(r.Body) > bodyLen {
				return 0, errors.New("content-length and acutal bosy length dosen't match")
			}
			break outer
		case StateDone:
			break outer
		}
	}

	return read, nil
}

var ErrBadRequestLine = errors.New("bad request line")
var ErrBadHTTPMethod = errors.New("method part is not capital")
var ErrBadHTTPVersion = fmt.Errorf("wrong http version")
var SEPERATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	headerLine := b[:idx]
	read := idx + len(SEPERATOR)
	headers := bytes.Split(headerLine, []byte(" "))
	if len(headers) != 3 {
		return nil, 0, ErrBadRequestLine
	}
	rl := &RequestLine{
		Method:        string(headers[0]),
		RequestTarget: string(headers[1]),
		HttpVersion:   string(headers[2]),
	}
	if !rl.ValidateMethod(string(headers[0])) {
		return nil, 0, ErrBadHTTPMethod
	}
	if !rl.ValidateHTTP() {
		return nil, 0, ErrBadHTTPVersion
	}
	rl.HttpVersion = strings.Split(rl.HttpVersion, "/")[1]
	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{State: StateInit, Headers: headers.NewHeaders()}
	buf := make([]byte, 1024)
	bufIdx := 0
	for !req.Done() {
		fmt.Println(req.State)
		n, err := reader.Read(buf[bufIdx:])
		if err != nil {
			if err == io.EOF {
				if !req.Done() {
					return nil, errors.New("unexpected EOF: body shorter than Content-Length")
				}
				break
			}
			return nil, err
		}
		bufIdx += n
		readN, err := req.parse(buf[:bufIdx])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufIdx])
		bufIdx -= readN
	}
	return req, nil
}
