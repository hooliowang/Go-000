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

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT)

	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-c:
				close(stop)
				return
			}
		}
	}()

	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return serveApp1(ctx, stop)
	})

	eg.Go(func() error {
		return serveApp2(ctx, stop)
	})

	eg.Go(func() error {
		return serverApp3(ctx, stop)
	})

	if err := eg.Wait(); err != nil {
		fmt.Println("err: ", err.Error())
		return
	}

	fmt.Println("exit")
}

func serveApp1(ctx context.Context, stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/r1", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("r1")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	defer func() {
		fmt.Println("server1 exit")
	}()

	go func() {
		<-stop
		fmt.Println("server1 shutdown")
		server.Shutdown(ctx)
	}()

	fmt.Println("start listen 8080")
	return server.ListenAndServe()
}

func serveApp2(ctx context.Context, stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/r2", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("r2")
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	defer func() {
		fmt.Println("server2 exit")
	}()

	go func() {
		<-stop
		fmt.Println("server2 shutdown")
		server.Shutdown(ctx)
	}()

	fmt.Println("start listen 8081")
	return server.ListenAndServe()
}

func serverApp3(ctx context.Context, stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/r3", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("r3")
	})

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	defer func() {
		fmt.Println("server3 exit")
	}()

	go func() {
		<-stop
		fmt.Println("server3 shutdown")
		server.Shutdown(ctx)
	}()

	fmt.Println("start listen 8082")
	return server.ListenAndServe()
}
