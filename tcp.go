package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const MAX_BYTES = 1024
const END_CONNECTION = "END"

func main() {
	f, err := setupLogging()
	if err != nil {
		return
	}
	defer f.Close()

	ln, err := setupServer()
	if err != nil {
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	log.Println("Started connection")
	defer func() {
		conn.Close()
		log.Println("Closed connection")
	}()

	var buffer string
	data := make([]byte, MAX_BYTES)
	bytesProcessed := 0
	for {
		i, err := conn.Read(data)

		if err != nil {
			if err == io.EOF {
				log.Println(buffer)
				return
			} else {
				log.Fatal("Could not read client data", err)
			}
		}

		bytesProcessed += i
		if bytesProcessed >= MAX_BYTES {
			log.Fatal("Client sent too much data", buffer)
			panic("Too much data")
		}
		buffer = buffer + string(data)
	}
}

func setupServer() (net.Listener, error) {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	fmt.Println("Starting a TCP server...")

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Could not listen on port ", port, err)
	}

	fmt.Printf("Now listening on on port %v\n\n", port)
	return ln, err
}

func setupLogging() (*os.File, error) {
	f, err := os.OpenFile("development.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f, err
}
