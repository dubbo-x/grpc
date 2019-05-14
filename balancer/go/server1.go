package main

import (
	"golang.org/x/net/context"
	"go/helloworld"
	"net"
	"log"
	"google.golang.org/grpc"
)

var addr1 = "localhost:50051"

type greeter1 struct {}

func (s *greeter1) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf(addr1)
	return &helloworld.HelloReply{Message: request.Name}, nil
}

func main() {
	listen, err := net.Listen("tcp", addr1)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &greeter1{})
	s.Serve(listen)
}