package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/changzhijay2/geekbang_go/week04/internal/server"
	"github.com/changzhijay2/geekbang_go/week04/pkg/config"
	"golang.org/x/sync/errgroup"
)

type Program struct {
	server *server.Server
	ctx    context.Context
	cancel context.CancelFunc
}

func NewProgram() *Program {
	server, err := InitServer()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Program{
		server: server,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Program) Run() error {
	g, ctx := errgroup.WithContext(p.ctx)
	g.Go(func() error {
		return p.server.ListenAndServer()
	})
	g.Go(func() error {
		<-ctx.Done()
		return p.server.Stop()
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-quit:
				p.Stop()
			}
		}
	})

	return g.Wait()
}

func (p *Program) Stop() error {
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "", "conf file path.")
	flag.Parse()
	config.InitConf(confPath)
	program := NewProgram()
	if err := program.Run(); err != nil {
		if !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}
}
