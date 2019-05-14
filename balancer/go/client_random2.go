package main

import (
	"google.golang.org/grpc/naming"
	"time"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"golang.org/x/net/context"
	"go/helloworld"
	"sync"
)

// watcher
type SimpleWatcher struct {
	index int
	addrs []string
}

func NewSimpleWatcher(addrs ...string) naming.Watcher {
	return &SimpleWatcher{index: 0, addrs: addrs}
}

func (w *SimpleWatcher) Next() ([]*naming.Update, error) {
	time.Sleep(2 * time.Second)

	addrs := make([]*naming.Update, 0, 2)
	for index := 0; index < 2; index++ {
		addr := w.addrs[w.index % 3]
		addrs = append(addrs, &naming.Update{Addr: addr})
		w.index++
	}

	return addrs, nil
}

func (w *SimpleWatcher) Close() {
	return
}

type SimpleResolver struct {}

// resolver
func NewSimpleResolver() naming.Resolver {
	return &SimpleResolver{}
}

func (r *SimpleResolver) Resolve(target string) (naming.Watcher, error) {
	return NewSimpleWatcher("localhost:50051", "localhost:50052", "localhost:50053"), nil
}

// balancer
type RandomBalancer2 struct {
	resolver naming.Resolver
	watcher naming.Watcher

	rw sync.RWMutex
	addrs []string
	addrsChan chan []grpc.Address

	wait *sync.Cond
	ready bool

	done bool
}

func NewRandomBalancer2(resolver naming.Resolver) grpc.Balancer {
	 balancer := &RandomBalancer2{
		resolver: resolver,
		addrsChan: make(chan []grpc.Address, 1),
		wait: sync.NewCond(&sync.Mutex{}),
	}
	return balancer
}

func (b *RandomBalancer2) Start(target string, config grpc.BalancerConfig) error {
	log.Println("Start", target)
	watcher, err := b.resolver.Resolve(target)
	if err != nil {
		return err
	}
	b.watcher = watcher

	go func() {
		for !b.done {
			updates, err := b.watcher.Next()
			if err != nil {
				continue
			}

			addrs := make([]string, 0)
			grpcAddrs := make([]grpc.Address, 0)
			for _, update := range updates {
				addrs = append(addrs, update.Addr)
				grpcAddrs = append(grpcAddrs, grpc.Address{Addr: update.Addr})
			}

			b.rw.Lock()
			b.addrsChan <- grpcAddrs
			b.addrs = addrs
			b.rw.Unlock()

			b.wait.L.Lock()
			b.ready = true
			b.wait.Broadcast()
			b.wait.L.Unlock()
		}
	}()

	return nil
}

func (b *RandomBalancer2) Up(addr grpc.Address) func(error) {
	log.Println("Up", addr)
	return nil
}

func (b *RandomBalancer2) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	b.wait.L.Lock()
	if !b.ready {
		b.wait.Wait()
	}
	b.wait.L.Unlock()

	b.rw.RLock()
	defer b.rw.RUnlock()
	length := len(b.addrs)
	index := rand.Intn(length)
	log.Println("Get", b.addrs[index])
	return grpc.Address{Addr: b.addrs[index]}, nil, nil
}

func (b *RandomBalancer2) Notify() <-chan []grpc.Address {
	return b.addrsChan
}

func (b *RandomBalancer2) Close() error {
	b.done = false
	return nil
}

func main() {
	r := NewSimpleResolver()
	b := NewRandomBalancer2(r)
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
		time.Sleep(time.Second)
	}
}