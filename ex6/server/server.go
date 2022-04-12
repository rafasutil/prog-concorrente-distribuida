package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

type Server int

var b bytes.Buffer

func (s *Server) Compress(req []byte, reply *[]byte) error {
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write(req); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}
	*reply = b.Bytes()
	return nil
}

/* --- Servidor --- */
func servidor() {

	fileCompress := new(Server)

	server := rpc.NewServer()
	err := server.RegisterName("File", fileCompress)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	ln, err := net.Listen("tcp", "localhost:1313")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	defer func(ln net.Listener) {
		var err = ln.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}(ln)

	fmt.Println("Servidor está pronto...")

	server.Accept(ln)

}

func main() {
	go servidor()
	_, _ = fmt.Scanln()
}
