package lab0

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type listenerInterface func(string, int, handlerInterface)

type handlerInterface func(conn net.Conn)

func TCPListener(host string, port int, handler handlerInterface) {
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
		go handler(conn)
	}
}

func TCPHandler(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	log.Printf("Connection from %s", clientAddr)

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Connection from %s closed", clientAddr)
				return
			}
			log.Printf("Error while reading data from %s: %s", clientAddr, err)
			return
		}

		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Printf("Error while sending data to %s: %s", clientAddr, err)
			return
		}
	}
}
