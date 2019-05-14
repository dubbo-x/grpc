package main

import (
	"golang.org/x/net/context"
	"go/helloworld"
	"net"
	"log"
	"google.golang.org/grpc"
)

var addr3 = "localhost:50053"

type greeter3 struct {}

func (s *greeter3) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Print(addr3)
	return &helloworld.HelloReply{Message: request.Name}, nil
}

func main() {
	listen, err := net.Listen("tcp", addr3)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &greeter3{})
	s.Serve(listen)
}