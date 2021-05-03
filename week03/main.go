package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/sync/errgroup"
)

type HttpServer struct {
	ctx     context.Context
	host    string
	port    int
	handler *httprouter.Router
	s       *http.Server
}

func NewHttpServer(host string, port int) *HttpServer {
	handler := httprouter.New()
	addr := fmt.Sprintf("%s:%d", host, port)
	return &HttpServer{
		ctx:     context.Background(),
		host:    host,
		port:    port,
		handler: handler,
		s: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (hs *HttpServer) RegisterRoute() {
	hs.handler.GET("/ping", hs.ping)
	hs.handler.GET("/stop", hs.stop)
}

func (hs *HttpServer) ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("pong..."))
}

func (hs *HttpServer) stop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("bye..."))
	hs.Shutdown()
}

func (hs *HttpServer) ListenAndServe() error {
	addr := fmt.Sprintf("%s:%d", hs.host, hs.port)
	log.Printf("http server listen on %s\n", addr)
	return hs.s.ListenAndServe()
}

func (hs *HttpServer) Run() error {
	hs.RegisterRoute()
	g, ctx := errgroup.WithContext(hs.ctx)

	g.Go(func() error {
		return hs.ListenAndServe()
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-quit:
				return hs.Shutdown()
			}
		}
	})

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (hs *HttpServer) Shutdown() error {
	log.Println("http server shutdown.")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return hs.s.Shutdown(ctx)
}

func main() {
	httpserver := NewHttpServer("0.0.0.0", 8081)
	if err := httpserver.Run(); err != nil {
		log.Printf("%+v", err)
		return
	}
}
