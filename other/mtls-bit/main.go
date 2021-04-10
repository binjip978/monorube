package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type cfg struct {
	waitTime time.Duration
	addr     string
	peerAddr string
	id int
}

func parseConf() cfg {
	c := cfg{}
	flag.StringVar(&c.addr, "addr", "localhost:8081", "server addr")
	flag.StringVar(&c.peerAddr, "peerAddr", "localhost:8081", "peer addr")
	flag.DurationVar(&c.waitTime, "dur", 2 * time.Second, "dur")
	flag.IntVar(&c.id, "id", 1, "id")
	flag.Parse()

	return c
}

func main() {
	cfg := parseConf()
	log.Println(cfg)

	var myCertPath string
	var peerCertPath string

	if cfg.id == 1 {
		myCertPath = "certs/1/"
		peerCertPath = "certs/2/"
	} else {
		myCertPath = "certs/2/"
		peerCertPath = "certs/1/"
	}

	go func() {
		// pin server certificate
		cert, err := ioutil.ReadFile(peerCertPath + "cert.pem")
		if err != nil {
			panic(err)
		}

		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(cert); !ok {
			panic(err)
		}

		tp := http.Transport{
			TLSClientConfig: &tls.Config{
				CurvePreferences: []tls.CurveID{tls.CurveP256},
				MinVersion: tls.VersionTLS12,
				//InsecureSkipVerify: true,
				RootCAs: certPool,
			},
		}

		client := http.Client{
			Timeout: 1 * time.Second,
			Transport: &tp,
		}

		for {
			<- time.After(cfg.waitTime)
			resp, err := client.Get(fmt.Sprintf("https://%s/bit", cfg.peerAddr))
			if err != nil {
				log.Println(err)
				continue
			}

			_, _ = io.Copy(ioutil.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/bit", func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()
		_, _ = io.Copy(ioutil.Discard, r.Body)
		log.Println(r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	})

	srv := http.Server{
		Addr: cfg.addr,
		Handler: mux,
		ReadHeaderTimeout: time.Second,
		IdleTimeout: 5 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS(myCertPath + "cert.pem", myCertPath + "key.pem"))
}
