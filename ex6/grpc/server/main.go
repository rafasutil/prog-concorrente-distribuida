package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net"
	"google.golang.org/grpc"
	pb "ex6/fileservice"
	"log"
)

type server struct{
	pb.UnimplementedFileServiceServer
}

var b bytes.Buffer

func (s *server) Compress(ctx context.Context, file *pb.FileRequest) (*pb.FileResponse, error) {
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write(file.File); err != nil {
		log.Fatal(err)
		return &pb.FileResponse{File: []byte{}}, err
	}

	if err := gz.Close(); err != nil {
		log.Fatal(err)
		return &pb.FileResponse{File: []byte{}}, err
	}

	return &pb.FileResponse{File: b.Bytes()}, nil
}

func servidor(){
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err)
	}
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})
	if err:= s.Serve(lis); err != nil {
		fmt.Println(err)
	}
}

func main() {
	go servidor()
	_, _ = fmt.Scanln()
}