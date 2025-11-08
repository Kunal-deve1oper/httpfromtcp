package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	remoteaddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("error occured: ", err)
	}
	localaddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		log.Fatal("error occured: ", err)
	}
	conn, err := net.DialUDP("udp", localaddr, remoteaddr)
	if err != nil {
		log.Fatal("error occured: ", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Println(">")
		text, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}
