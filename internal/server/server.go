package server

import (
	"fmt"
	"io"
	"net"

	"github.com/Kunal-deve1oper/httpfromtcp/internal/request"
	"github.com/Kunal-deve1oper/httpfromtcp/internal/response"
)

type Server struct {
	connectionClose bool
	handler         Handler
}

type HandlerError struct {
	statusCode response.StatusCode
	message    string
}

func (h *HandlerError) WriteError(status response.StatusCode, message string) {
	h.statusCode = status
	h.message = message
}

type Handler func(w *response.Writer, req *request.Request) *HandlerError

func (s *Server) Close() {
	s.connectionClose = true
}

func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()
	writer := &response.Writer{Data: conn}
	h := response.GetDefaultHeaders(0)
	request, err := request.RequestFromReader(conn)
	if err != nil {
		writer.WriteStatusLine(response.StatusBadRequest)
		writer.WriteHeaders(h)
		return
	}
	handlerErr := s.handler(writer, request)
	if handlerErr != nil {
		writer.WriteStatusLine(handlerErr.statusCode)
		body := []byte(handlerErr.message)
		h.Set("content-length", fmt.Sprintf("%d", (len(body))))
		h.Set("content-type", "text/html")
		writer.WriteHeaders(h)
		conn.Write(body)
	}
}

func (s *Server) runServer(listner net.Listener) {
	for {
		conn, err := listner.Accept()
		if s.connectionClose {
			return
		}
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}

func Serve(port int, handler Handler) (*Server, error) {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{connectionClose: false, handler: handler}
	go server.runServer(listner)
	return server, nil
}
