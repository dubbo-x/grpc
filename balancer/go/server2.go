package main

import (
	"golang.org/x/net/context"
	"go/helloworld"
	"net"
	"log"
	"google.golang.org/grpc"
)

var addr2 = "localhost:50052"

type greeter2 struct {}

func (s *greeter2) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Print(addr2)
	return &helloworld.HelloReply{Message: request.Name}, nil
}

func main() {
	listen, err := net.Listen("tcp", addr2)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &greeter2{})
	s.Serve(listen)
}