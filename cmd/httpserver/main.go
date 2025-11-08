package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Kunal-deve1oper/httpfromtcp/internal/request"
	"github.com/Kunal-deve1oper/httpfromtcp/internal/response"
	"github.com/Kunal-deve1oper/httpfromtcp/internal/server"
)

const port = 42069

func handler(w *response.Writer, r *request.Request) *server.HandlerError {
	handlerError := &server.HandlerError{}
	h := response.GetDefaultHeaders(0)
	reqLine := r.RequestLine.RequestTarget
	if reqLine == "/yourproblem" {
		body := `<html>
					<head>
						<title>400 Bad Request</title>
					</head>
					<body>
						<h1>Bad Request</h1>
						<p>Your request honestly kinda sucked.</p>
					</body>
					</html>`
		handlerError.WriteError(response.StatusBadRequest, body)
		return handlerError
	} else if reqLine == "/myproblem" {
		body := `<html>
					<head>
						<title>500 Internal Server Error</title>
					</head>
					<body>
						<h1>Internal Server Error</h1>
						<p>Okay, you know what? This one is on me.</p>
					</body>
					</html>`
		handlerError.WriteError(response.StatusInternalError, body)
		return handlerError
	} else if strings.HasPrefix("/httpbin/stream", reqLine) {
		res, err := http.Get("https://httpbin.org/" + reqLine[len("/httpbin/"):])
		if err != nil {
			body := `<html>
					<head>
						<title>500 Internal Server Error</title>
					</head>
					<body>
						<h1>Internal Server Error</h1>
						<p>Okay, you know what? This one is on me.</p>
					</body>
					</html>`
			handlerError.WriteError(response.StatusInternalError, body)
			return handlerError
		}
		w.WriteStatusLine(response.StatusOk)
		h.Set("content-length", "text/plain")
		h.Set("Transfer-Encoding", "chunked")
		for {
			buff := make([]byte, 32)
			n, err := res.Body.Read(buff)
			if err != nil {
				break
			}
			w.WriteChunkedBody(buff[:n])
		}
		w.WriteChunkedBodyDone()
	} else {
		body := `<html>
					<head>
						<title>200 OK</title>
					</head>
					<body>
						<h1>Success!</h1>
						<p>Your request was an absolute banger.</p>
					</body>
					</html>`
		w.WriteStatusLine(response.StatusOk)
		h.Set("content-length", fmt.Sprintf("%d", len(body)))
		h.Set("content-type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
	}
	return nil
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
