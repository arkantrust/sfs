package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"bytes"
	"encoding/binary"
)

type FileServer struct{}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go fs.handleConnection(conn)
	}
}

func (fs *FileServer) handleConnection(conn net.Conn) {
	buf := new(bytes.Buffer)
	defer conn.Close()
	for {
        var size int64
        binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed the connection")
			} else {
				fmt.Println("Error reading from connection:", err)
			}
			return
		}
        if n > 0 {
            fmt.Println(buf.Bytes()) // TODO: Convert to string for better readability
            fmt.Printf("received %d bytes\n", n)
        }
	}
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}
	defer conn.Close()

    binary.Write(conn, binary.LittleEndian, int64(size))
    n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("sent %d bytes\n", n)
	return nil
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		if err := sendFile(30000); err != nil {
			log.Println("Error sending file:", err)
		}
	}()
	server := &FileServer{}
	server.start()
}