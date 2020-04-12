package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
)

type args struct {
	address  string
	begin    int
	end      int
	poolSize int
}

func parseArgs() args {
	address := flag.String("addr", "localhost", "address")
	begin := flag.Int("begin", 1, "start of port range")
	end := flag.Int("end", 1024, "end of port range")
	poolSize := flag.Int("size", 32, "workers pool size")
	flag.Parse()

	return args{
		address:  *address,
		begin:    *begin,
		end:      *end,
		poolSize: *poolSize,
	}
}

type results struct {
	lock      sync.Mutex
	openPorts []int
}

func main() {
	args := parseArgs()
	ports := make(chan int, args.poolSize*3)
	defer close(ports)
	var wg sync.WaitGroup
	results := &results{}

	for i := 0; i < args.poolSize; i++ {
		go scan(ports, results, &wg, args.address)
	}

	for i := args.begin; i <= args.end; i++ {
		wg.Add(1)
		ports <- i
	}

	wg.Wait()

	results.lock.Lock()
	defer results.lock.Unlock()
	sort.Ints(results.openPorts)
	for _, openPort := range results.openPorts {
		fmt.Printf("%5d open\n", openPort)
	}
}

func scan(ports chan int, results *results, wg *sync.WaitGroup, addr string) {
	for port := range ports {
		address := fmt.Sprintf("%s:%d", addr, port)
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			results.lock.Lock()
			results.openPorts = append(results.openPorts, port)
			results.lock.Unlock()
		}
		wg.Done()
	}
}
