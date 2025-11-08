package response

import (
	"fmt"
	"io"

	"github.com/Kunal-deve1oper/httpfromtcp/internal/headers"
)

type StatusCode string

type Writer struct {
	Data io.Writer
}

const (
	StatusOk            StatusCode = "200"
	StatusBadRequest    StatusCode = "400"
	StatusInternalError StatusCode = "500"
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusOk:
		_, err := w.Data.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case StatusBadRequest:
		_, err := w.Data.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case StatusInternalError:
		_, err := w.Data.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	default:
		w.Data.Write([]byte("unknown body\r\n"))
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	r := headers.Headers{}
	r.Set("content-length", fmt.Sprintf("%d", contentLen))
	r.Set("connection", "close")
	r.Set("content-type", "text/plain")
	return r
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	contentLength, ok := headers.Get("content-length")
	if !ok {
		contentLength = ""
	}
	connection, ok := headers.Get("connection")
	if !ok {
		connection = ""
	}
	contentType, ok := headers.Get("content-type")
	if !ok {
		contentType = ""
	}
	resp := fmt.Sprintf("Content-Length: %s\r\nConnection: %s\r\nContent-Type: %s\r\n\r\n", contentLength, connection, contentType)
	_, err := w.Data.Write([]byte(resp))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Data.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunk := []byte(fmt.Sprintf("%x\r\n", len(p)))
	chunk = append(chunk, p...)
	chunk = append(chunk, []byte("\r\n")...)
	return w.WriteBody(chunk)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.WriteBody([]byte("0\r\n\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	return w.WriteHeaders(h)
}
