package main

import (
	"bufio"
	"errors"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:5433")
	if err != nil {
		panic(err)
	}
	log.Printf("server is running...\n")

	for {
		if err := processPayment(listener); err != nil {
			log.Printf("client conn err: %s", err)
		}
	}
}

func processPayment(lis net.Listener) error {
	conn, err := lis.Accept()
	if err != nil {
		return errors.New("accept new conn ")
	}
	defer conn.Close()

	log.Printf("receive conn")
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return errors.New(" get err from client")
	}
	log.Printf("message received %s", msg)

	if _, err := conn.Write([]byte("ok\n")); err != nil {
		return errors.New("err write response to client")
	}
	return nil
}


