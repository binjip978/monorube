package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go.bug.st/serial.v1"
)

func readTemp(device string, rate int) (<-chan float64, error) {
	ch := make(chan float64)
	mode := &serial.Mode{
		BaudRate: rate,
	}
	port, err := serial.Open(device, mode)
	if err != nil {
		return ch, err
	}

	go func() {
		b := make([]byte, 16)
		var buff bytes.Buffer
		for {
			n, _ := port.Read(b)
			for i := 0; i < n; i++ {
				switch b[i] {
				case '\r': // skip
				case '\n':
					s := buff.String()
					buff.Reset()
					f, err := strconv.ParseFloat(s, 32)
					if err == nil {
						ch <- f
					}
				default:
					buff.WriteByte(b[i])
				}
			}
		}
	}()

	return ch, nil
}

type cfg struct {
	device     string
	rate       int
	server     bool
	serverPort string
}

func parseConfig() cfg {
	var device = flag.String("dev", "/dev/ttyACM0", "device")
	var port = flag.Int("port", 9600, "baud rate")
	var server = flag.Bool("server", false, "run http server, on /weather")
	var serverPort = flag.String("srvPort", "8080", "server port")
	flag.Parse()

	return cfg{*device, *port, *server, *serverPort}
}

type server struct {
	last float64
	srv  *http.Server
	sync.Mutex
}

func newServer(ch <-chan float64, port string) *server {
	srv := &server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/weather", func(w http.ResponseWriter, req *http.Request) {
		srv.Lock()
		defer srv.Unlock()
		w.Write([]byte(fmt.Sprintf("%.2f\n", srv.last)))
	})

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
	srv.srv = s

	return srv
}

func main() {
	cfg := parseConfig()
	ch, err := readTemp(cfg.device, cfg.rate)
	if err != nil {
		panic(err)
	}

	if cfg.server {
		srv := newServer(ch, cfg.serverPort)
		go func() {
			for v := range ch {
				srv.Lock()
				srv.last = v
				srv.Unlock()
			}
		}()

		log.Fatal(srv.srv.ListenAndServe())
	}

	// console mode
	for v := range ch {
		fmt.Printf("%.2f\n", v)
	}
}
