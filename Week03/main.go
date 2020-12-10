package main

import (
	"context"
	"fmt"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

var (
	stopChs = make(map[string]chan<- struct{})
)

type HttpServer struct {
	svr    *http.Server
	name   string
	stopCh chan struct{} // receive the quit signal
}

func NewServer(addr string) (server *HttpServer) {
	server = &HttpServer{}
	server.svr = &http.Server{
		Addr: addr,
	}
	server.stopCh = make(chan struct{})

	return
}

func (server *HttpServer) RegistRouter(uri string, handler func(resp http.ResponseWriter, req *http.Request)) {
	mux := http.NewServeMux()
	mux.HandleFunc(uri, handler)
	server.svr.Handler = mux
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT)

	stop := make(chan struct{}) // receive the ctrl+c signal and close all servers

	// ctrl+c
	go func() {
		for {
			select {
			case <-c:
				close(stop)
				return
			}
		}
	}()

	eg, _ := errgroup.WithContext(context.Background())
	// start server 1
	eg.Go(func() error {
		return serveApp1(stop)
	})

	// start server 2
	eg.Go(func() error {
		return serveApp2(stop)
	})

	// start server 3 for shutdown assigned server
	eg.Go(func() error {
		return serverApp3(stop)
	})

	if err := eg.Wait(); err != nil {
		fmt.Println("err: ", err.Error())
	}

	fmt.Println("exit")
}

func serveApp1(stop <-chan struct{}) error {
	server := NewServer(":8080")
	server.name = "s1"
	server.stopCh = make(chan struct{})
	server.RegistRouter("/r1", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("r1")
	})

	stopChs[server.name] = server.stopCh

	defer func() {
		fmt.Println("server1 exit")
	}()

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("server1 shutdown")
				server.svr.Shutdown(context.Background())
				return
			case <-server.stopCh:
				fmt.Println("server1 shutdown itself")
				server.svr.Shutdown(context.Background())
				return
			}
		}
	}()

	fmt.Println("start listen 8080")
	return server.svr.ListenAndServe()
}

func serveApp2(stop <-chan struct{}) error {
	server := NewServer(":8081")
	server.name = "s2"
	server.stopCh = make(chan struct{})
	server.RegistRouter("/r2", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("r2")
	})

	stopChs[server.name] = server.stopCh

	defer func() {
		fmt.Println("server2 exit")
	}()

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("server2 shutdown")
				server.svr.Shutdown(context.Background())
				return
			case <-server.stopCh:
				fmt.Println("server2 shutdown itself")
				server.svr.Shutdown(context.Background())
				return
			}
		}
	}()

	fmt.Println("start listen 8081")
	return server.svr.ListenAndServe()
}

func serverApp3(stop chan struct{}) error {
	server := NewServer(":8082")
	server.name = "s3"
	server.stopCh = make(chan struct{})
	server.RegistRouter("/close", func(resp http.ResponseWriter, req *http.Request) {
		servername := req.FormValue("servername")
		fmt.Println("servername: ", servername)

		if servername == "all" {
			close(stop)
		} else {
			ch, exist := stopChs[servername]
			if exist {
				close(ch)
				delete(stopChs, servername)
			}
		}
	})

	defer func() {
		fmt.Println("server3 exit")
	}()

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("server3 shutdown")
				server.svr.Shutdown(context.Background())
				return
			}
		}
	}()

	fmt.Println("start listen 8082")
	return server.svr.ListenAndServe()
}
