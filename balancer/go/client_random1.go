package main

import (
	"log"
	"go/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
)

// balancer
type RandomBalancer struct {
	addrs []string
	addrsChan chan []grpc.Address
}

func NewRandomBalancer(addrs...string) grpc.Balancer {
	grpcAddrs := make([]grpc.Address, 0, len(addrs))
	for _, addr := range addrs {
		grpcAddrs = append(grpcAddrs, grpc.Address{Addr: addr})
	}

	addrsChan := make(chan []grpc.Address, 1)
	addrsChan <- grpcAddrs

	return &RandomBalancer{addrs: addrs, addrsChan: addrsChan}
}

func (b *RandomBalancer) Start(target string, config grpc.BalancerConfig) error {
	log.Println("Start", target)
	return nil
}

func (b *RandomBalancer) Up(addr grpc.Address) func(error) {
	log.Println("Up", addr)
	return nil
}

func (b *RandomBalancer) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	length := len(b.addrs)
	index := rand.Intn(length)
	log.Println("Get", b.addrs[index])
	return grpc.Address{Addr: b.addrs[index]}, nil, nil
}

func (b *RandomBalancer) Notify() <-chan []grpc.Address {
	log.Println("Notify", b.addrsChan)
	return b.addrsChan
}

func (b *RandomBalancer) Close() error {
	log.Println("Close")
	return nil
}

func main() {
	b := NewRandomBalancer("localhost:50051", "localhost:50052", "localhost:50053")
	conn, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	for index := 0; index < 10; index++ {
		response, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "hello world"})
		if err != nil {
			log.Fatal(err)
		}
		log.Print(response.Message)
	}
}