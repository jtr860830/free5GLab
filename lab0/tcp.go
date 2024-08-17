package lab0

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type TCPListener interface {
	listen(string, int, TCPHandler)
}

type TCPHandler interface {
	handle(net.Conn)
}

type (
	TCPListenFunc func(string, int, TCPHandler)
	TCPHandleFunc func(conn net.Conn)
)

func (f TCPListenFunc) listen(host string, port int, handler TCPHandler) {
	f(host, port, handler)
}

func (f TCPHandleFunc) handle(conn net.Conn) {
	f(conn)
}

func Listen(host string, port int, handler TCPHandler) {
	srv, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Error while starting TCP server: %s", err)
	}
	defer srv.Close()
	log.Printf("TCP server started on %s:%d", host, port)

	for {
		conn, err := srv.Accept()
		if err != nil {
			log.Fatalf("Error while accepting connection: %s", err)
		}
		go handler.handle(conn)
	}
}

func Handle(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	log.Printf("Connection from %s", clientAddr)
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Connection from %s closed", clientAddr)
				return
			}
			log.Printf("Error while reading data from %s: %s", clientAddr, err)
			return
		}

		if _, err = conn.Write([]byte(data)); err != nil {
			log.Printf("Error while sending data to %s: %s", clientAddr, err)
			return
		}
	}
}
