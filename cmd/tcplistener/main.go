package main

import (
	"log"
	"net"

	"github.com/Kunal-deve1oper/httpfromtcp/internal/request"
)

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	data := make(chan string)
// 	buff := make([]byte, 8)
// 	lines := ""
// 	go func() {
// 		defer close(data)
// 		defer f.Close()
// 		defer func() {
// 			if len(lines) != 0 {
// 				data <- lines
// 			}
// 		}()
// 		for {
// 			n, err := f.Read(buff)
// 			chunk := buff[:n]
// 			if err == io.EOF {
// 				if len(lines) > 0 {
// 					data <- lines
// 				}
// 				return
// 			} else if err != nil {
// 				return
// 			}
// 			if i := bytes.IndexByte(chunk, '\n'); i != -1 {
// 				lines += string(chunk[:i])
// 				data <- lines
// 				lines = ""
// 				lines += string(chunk[i+1:])
// 			} else {
// 				lines += string(chunk)
// 			}
// 		}
// 	}()
// 	return data
// }

func main() {
	// file, err := os.Open("message.txt")
	// if err != nil {
	// 	log.Fatal("Cannot open file :", err)
	// }
	// defer file.Close()
	// data := getLinesChannel(file)
	// for line := range data {
	// 	fmt.Printf("read: %s\n", line)
	// }
	listner, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error occured")
	}
	defer listner.Close()
	for {
		conn, err := listner.Accept()
		log.Println("a new connecction has been opened: ", conn.LocalAddr().String())
		if err != nil {
			log.Fatal("error occured")
		}
		rl, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error occured from read line: ", err)
		}
		log.Println("Request line:")
		log.Printf("- Method: %s", rl.RequestLine.Method)
		log.Printf("- Target: %s", rl.RequestLine.RequestTarget)
		log.Printf("- Version: %s", rl.RequestLine.HttpVersion)
		log.Print("Headers:\n")
		for key, value := range rl.Headers {
			log.Printf("- %s: %s\n", key, value)
		}
		log.Println("Body:")
		log.Printf("%s\n", string(rl.Body))
		log.Println("connecction has been closed: ", conn.LocalAddr().String())
	}
}
